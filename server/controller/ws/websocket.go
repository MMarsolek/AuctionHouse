package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/MMarsolek/AuctionHouse/log"
	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/MMarsolek/AuctionHouse/server/controller/auth"
	"github.com/MMarsolek/AuctionHouse/server/controller/middleware"
	"github.com/MMarsolek/AuctionHouse/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func newErrorMessage(command SocketCommand, statusCode int, message string, args ...interface{}) *responseMessage {
	return &responseMessage{
		Command:    command,
		Message:    fmt.Sprintf(message, args...),
		StatusCode: statusCode,
	}
}

type sessionData struct {
	ws         *websocket.Conn
	ctx        context.Context
	username   string
	permission model.PermissionLevel
}

// Handler handles websocket connections and allows for clients to be updated when a bid is placed.
type Handler struct {
	upgrader           *websocket.Upgrader
	currentConnections map[string]*sessionData
	rwLock             sync.RWMutex

	userClient storage.UserClient
	itemClient storage.AuctionItemClient
	bidClient  storage.AuctionBidClient
}

// NewHandler constructs a Handler.
func NewHandler(
	userClient storage.UserClient,
	itemClient storage.AuctionItemClient,
	bidClient storage.AuctionBidClient,
) *Handler {
	return &Handler{
		upgrader: &websocket.Upgrader{
			HandshakeTimeout: time.Second * 30,
			Error:            handlerWSError,
		},
		currentConnections: make(map[string]*sessionData),
		userClient:         userClient,
		itemClient:         itemClient,
		bidClient:          bidClient,
	}
}

// RegisterRoutes registers the routes to establish a websocket connection.
func (handler *Handler) RegisterRoutes(rootRouter *mux.Router) {
	wsRouter := rootRouter.PathPrefix("/ws").Subrouter()
	wsRouter.Use(middleware.VerifyAuthToken)

	wsRouter.HandleFunc("", handler.ServeWS).Methods(http.MethodGet)
}

// ServeWS upgrades the connection to a websocket and listens for different JSON commands.
func (handler *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	ws, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(r.Context(), "unable to upgrade websocket", "err", err)
	}
	defer ws.Close()
	r = r.WithContext(log.WithFields(r.Context(), "client", getIdentifierFromWebsocket(ws)))

	func() {
		handler.rwLock.Lock()
		defer handler.rwLock.Unlock()

		handler.currentConnections[getIdentifierFromWebsocket(ws)] = &sessionData{
			ws:         ws,
			ctx:        r.Context(),
			username:   auth.ExtractUsername(r.Context()),
			permission: auth.ExtractPermission(r.Context()),
		}
	}()

	running := true
	for running {
		var message commandMessage
		err = ws.ReadJSON(&message)
		if err != nil {
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				func() {
					handler.rwLock.Lock()
					defer handler.rwLock.Unlock()
					running = false

					delete(handler.currentConnections, getIdentifierFromWebsocket(ws))
				}()
			}
			log.Error(r.Context(), "unable to read JSON from client", "err", err)
			ws.WriteJSON(newErrorMessage(SocketCommandUnknown, http.StatusBadRequest, "invalid request format"))
			continue
		}

		if message.Command == SocketCommandPlaceBid {
			err = handler.handlePlaceBid(handler.getSessionData(ws), message.Payload.(*commandMessagePlaceBid))
		} else {
			log.Error(r.Context(), "unrecognized command", "command", message.Command)
			continue
		}

		if err != nil {
			ws.WriteJSON(newErrorMessage(message.Command, http.StatusInternalServerError, "unhandled error: %v", err))
			log.Error(r.Context(), "unhandled error", "command", message.Command, "err", err)
		}
	}
}

// handlePlaceBid handles incoming commands where the user wants to place a bid.
func (handler *Handler) handlePlaceBid(data *sessionData, command *commandMessagePlaceBid) error {
	user, err := handler.userClient.Get(data.ctx, data.username)
	if err != nil {
		return errors.Wrap(err, "could not retrieve user")
	}

	item, err := handler.itemClient.Get(data.ctx, command.ItemName)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			if writeErr := data.ws.WriteJSON(newErrorMessage(
				SocketCommandPlaceBid,
				http.StatusNotFound,
				"could not find item '%s'",
				command.ItemName,
			)); writeErr != nil {
				return errors.Wrap(writeErr, "unable to write to client")
			}
			return nil
		}
		return errors.Wrap(err, "could not retrieve item")
	}

	_, err = handler.bidClient.PlaceBid(data.ctx, user, item, command.BidAmount)
	if err != nil {
		if errors.Is(err, storage.ErrBidTooLow) {
			if writeErr := data.ws.WriteJSON(newErrorMessage(
				SocketCommandPlaceBid,
				http.StatusBadRequest,
				"bid amount is too low",
			)); writeErr != nil {
				return errors.Wrap(writeErr, "unable to write to client")
			}
			return nil
		}
		return errors.Wrap(err, "unable to make bid")
	}

	func() {
		handler.rwLock.RLock()
		defer handler.rwLock.RUnlock()

		for _, connection := range handler.currentConnections {
			if err = connection.ws.WriteJSON(responseMessage{
				Command:    SocketCommandPlaceBid,
				StatusCode: http.StatusCreated,
				Message:    "New bid placed",
				Data: struct {
					ItemName string `json:"itemName"`
					Username string `json:"username"`
					NewBid   int    `json:"amount"`
				}{
					ItemName: item.Name,
					Username: user.Username,
					NewBid:   command.BidAmount,
				},
			}); err != nil {
				log.Error(connection.ctx, "unable to write to client", "command", SocketCommandPlaceBid)
			}
		}
	}()

	return nil
}

func (handler *Handler) getSessionData(ws *websocket.Conn) *sessionData {
	handler.rwLock.RLock()
	defer handler.rwLock.RUnlock()

	return handler.currentConnections[getIdentifierFromWebsocket(ws)]
}

func getIdentifierFromWebsocket(ws *websocket.Conn) string {
	return ws.RemoteAddr().String()
}

func handlerWSError(w http.ResponseWriter, r *http.Request, status int, errReason error) {
	w.WriteHeader(status)
	response, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("internal server websocket error: %v", errReason),
	})
	if err != nil {
		log.Error(r.Context(), "internal websocket error", "err", errReason)
		return
	}

	fmt.Fprintf(w, string(response))
}
