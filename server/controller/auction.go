package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MMarsolek/AuctionHouse/log"
	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/MMarsolek/AuctionHouse/server/controller/auth"
	"github.com/MMarsolek/AuctionHouse/server/controller/middleware"
	"github.com/MMarsolek/AuctionHouse/storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type (
	userResponse struct {
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
	}

	itemResponse struct {
		Name        string `json:"name"`
		ImageRef    string `json:"image,omitempty"`
		Description string `json:"description,omitempty"`
	}

	getHighestBidsResponse struct {
		BidAmount int           `json:"bidAmount"`
		Bidder    *userResponse `json:"bidder"`
		Item      *itemResponse `json:"item"`
	}

	getItemResponse struct {
		Name        string `json:"name"`
		ImageRef    string `json:"image,omitempty"`
		Description string `json:"description,omitempty"`
	}

	postItemRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ImageRef    string `json:"image"`
	}

	putItemRequest struct {
		Description string `json:"description,omitempty"`
		ImageRef    string `json:"image,omitempty"`
	}

	postBidRequest struct {
		BidAmount int `json:"amount"`
	}
)

// AuctionHandler provides handlers for endpoints involving model.AuctionItems and model.AuctionBids.
type AuctionHandler struct {
	userClient        storage.UserClient
	auctionItemClient storage.AuctionItemClient
	auctionBidClient  storage.AuctionBidClient
}

// NewAuctionHandler creates a new AuctionHandler with the necessary storage objects.
func NewAuctionHandler(
	userClient storage.UserClient,
	auctionItemClient storage.AuctionItemClient,
	auctionBidClient storage.AuctionBidClient,
) *AuctionHandler {
	return &AuctionHandler{
		userClient:        userClient,
		auctionItemClient: auctionItemClient,
		auctionBidClient:  auctionBidClient,
	}
}

// RegisterRoutes registers all of the paths to the handler functions.
func (handler *AuctionHandler) RegisterRoutes(router *mux.Router) {
	auctionsRouter := router.PathPrefix("/v1/auctions").Subrouter()
	auctionsRouter.Use(middleware.VerifyAuthToken)

	auctionsRouterAdmin := auctionsRouter.NewRoute().Subrouter()
	auctionsRouterAdmin.Use(middleware.VerifyPermissions(model.PermissionLevelAdmin))

	itemsAdmin := auctionsRouterAdmin.PathPrefix("/items").Subrouter()
	itemsAdmin.HandleFunc("", wrapHandler(handler.PostItem)).Methods(http.MethodPost)
	itemsAdmin.HandleFunc("/{itemName}", wrapHandler(handler.PutItem)).Methods(http.MethodPut)
	itemsAdmin.HandleFunc("/{itemName}", wrapHandler(handler.DeleteItem)).Methods(http.MethodDelete)

	auctionsRouterBidder := auctionsRouter.NewRoute().Subrouter()
	auctionsRouterBidder.Use(middleware.VerifyPermissions(model.PermissionLevelBidder))

	bidsBidder := auctionsRouterBidder.PathPrefix("/bids").Subrouter()
	bidsBidder.HandleFunc("/{itemName}", wrapHandler(handler.PostBid)).Methods(http.MethodPost)

	auctionsRouterBoth := auctionsRouter.NewRoute().Subrouter()
	auctionsRouterBoth.Use(middleware.VerifyPermissions(model.PermissionLevelAdmin, model.PermissionLevelBidder))

	itemsBoth := auctionsRouterBoth.PathPrefix("/items").Subrouter()
	itemsBoth.HandleFunc("", wrapHandler(handler.GetItems)).Methods(http.MethodGet)
	itemsBoth.HandleFunc("/{itemName}", wrapHandler(handler.GetItem)).Methods(http.MethodGet)

	bidsBoth := auctionsRouterBoth.PathPrefix("/bids").Subrouter()
	bidsBoth.HandleFunc("/{itemName}", wrapHandler(handler.GetHighestBid)).Methods(http.MethodGet)
	bidsBoth.HandleFunc("", wrapHandler(handler.GetHighestBids)).Methods(http.MethodGet)

}

// GetItem is the handler that retrieves the model.AuctionItem as serialized JSON.
func (handler *AuctionHandler) GetItem(w http.ResponseWriter, r *http.Request) error {
	itemName := mux.Vars(r)["itemName"]
	r = r.WithContext(log.WithFields(r.Context(), "itemName", itemName))
	item, err := handler.auctionItemClient.Get(r.Context(), itemName)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return nil
		}
		return errors.Wrap(err, "could not retrieve auction item")
	}

	rawItem, err := json.Marshal(getItemResponse{
		Name:        item.Name,
		ImageRef:    item.ImageRef,
		Description: item.Description,
	})

	if err != nil {
		return errors.Wrap(err, "could not marshal item")
	}

	fmt.Fprint(w, string(rawItem))
	return nil
}

// GetItems is the handler that retrieves all model.AuctionItems as serialized JSON.
func (handler *AuctionHandler) GetItems(w http.ResponseWriter, r *http.Request) error {
	items, err := handler.auctionItemClient.GetAll(r.Context())
	if err != nil {
		return errors.Wrap(err, "could not retrieve auction item")
	}

	responseObjects := make([]*getItemResponse, len(items))
	for i, item := range items {
		responseObjects[i] = &getItemResponse{
			Name:        item.Name,
			ImageRef:    item.ImageRef,
			Description: item.Description,
		}
	}

	rawResponse, err := json.Marshal(responseObjects)
	if err != nil {
		return errors.Wrap(err, "could not marshal item")
	}

	fmt.Fprint(w, string(rawResponse))
	return nil
}

// PostItem is the handler that creates a new model.AuctionItem.
func (handler *AuctionHandler) PostItem(w http.ResponseWriter, r *http.Request) error {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "could not read body")
	}

	var request postItemRequest
	err = json.Unmarshal(rawBody, &request)
	if err != nil {
		log.Info(r.Context(), "invalid json request", "body", string(rawBody), "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	defer r.Body.Close()

	newItem := &model.AuctionItem{
		Name:        request.Name,
		Description: request.Description,
		ImageRef:    request.ImageRef,
	}

	err = handler.auctionItemClient.Create(r.Context(), newItem)
	if err != nil {
		if errors.Is(err, storage.ErrEntityAlreadyExists) {
			w.WriteHeader(http.StatusBadRequest)
			response, marshalErr := json.Marshal(newErrorResponse("item already exists"))
			if marshalErr != nil {
				return errors.Wrap(marshalErr, "could not marshal error response")
			}
			fmt.Fprint(w, string(response))
			return nil
		}
		return errors.Wrap(err, "could not store item")
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

// PutItem is the handler that updates the fields of a model.AuctionItem.
func (handler *AuctionHandler) PutItem(w http.ResponseWriter, r *http.Request) error {
	itemName := mux.Vars(r)["itemName"]

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "could not read body")
	}

	var request putItemRequest
	err = json.Unmarshal(rawBody, &request)
	if err != nil {
		log.Info(r.Context(), "invalid json request", "body", string(rawBody), "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	defer r.Body.Close()

	updateFields := &model.AuctionItem{
		Name:        itemName,
		Description: request.Description,
		ImageRef:    request.ImageRef,
	}

	err = handler.auctionItemClient.Update(r.Context(), updateFields)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
			response, marshalErr := json.Marshal(newErrorResponse("item does not exist"))
			if marshalErr != nil {
				return errors.Wrap(marshalErr, "could not marshal error response")
			}
			fmt.Fprint(w, string(response))
			return nil
		}
		return errors.Wrap(err, "could not store item")
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

// DeleteItem is the handler that removes a model.AuctionItem from storage.
func (handler *AuctionHandler) DeleteItem(w http.ResponseWriter, r *http.Request) error {
	itemName := mux.Vars(r)["itemName"]

	err := handler.auctionItemClient.Delete(r.Context(), itemName)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
			response, marshalErr := json.Marshal(newErrorResponse("item does not exist"))
			if marshalErr != nil {
				return errors.Wrap(marshalErr, "could not marshal error response")
			}
			fmt.Fprint(w, string(response))
			return nil
		}
		return errors.Wrap(err, "could not delete item")
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

// GetHighestBid is the handler that finds the highest bid for a model.AuctionItem and retrieves the model.AuctionBid as
// serialized JSON.
func (handler *AuctionHandler) GetHighestBid(w http.ResponseWriter, r *http.Request) error {
	itemName := mux.Vars(r)["itemName"]

	item, err := handler.auctionItemClient.Get(r.Context(), itemName)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
			response, marshalErr := json.Marshal(newErrorResponse("item does not exist"))
			if marshalErr != nil {
				return errors.Wrap(marshalErr, "could not marshal error response")
			}
			fmt.Fprint(w, string(response))
			return nil
		}
		return errors.Wrap(err, "could not get item")
	}

	highestBid, err := handler.auctionBidClient.GetHighestBid(r.Context(), item)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
			response, marshalErr := json.Marshal(newErrorResponse("no bids for item"))
			if marshalErr != nil {
				return errors.Wrap(marshalErr, "could not marshal error response")
			}
			fmt.Fprint(w, string(response))
			return nil
		}
		return errors.Wrap(err, "could not get item")
	}

	rawResponse, err := json.Marshal(&getHighestBidsResponse{
		BidAmount: highestBid.BidAmount,
		Bidder: &userResponse{
			Username:    highestBid.Bidder.Username,
			DisplayName: highestBid.Bidder.DisplayName,
		},
		Item: &itemResponse{
			Name:        highestBid.Item.Name,
			Description: highestBid.Item.Description,
			ImageRef:    highestBid.Item.ImageRef,
		},
	})
	if err != nil {
		errors.Wrap(err, "unable to marshal highest bid response")
	}

	fmt.Fprintf(w, string(rawResponse))
	return nil
}

// GetHighestBid is the handler that finds the highest bids all model.AuctionItems and retrieves all model.AuctionBids as
// serialized JSON.
func (handler *AuctionHandler) GetHighestBids(w http.ResponseWriter, r *http.Request) error {
	highestBids, err := handler.auctionBidClient.GetAllHighestBids(r.Context())
	if err != nil {
		return errors.Wrap(err, "unable to get highest bids")
	}

	responseObjects := make([]*getHighestBidsResponse, len(highestBids))
	for i, highestBid := range highestBids {
		responseObjects[i] = &getHighestBidsResponse{
			BidAmount: highestBid.BidAmount,
			Bidder: &userResponse{
				Username:    highestBid.Bidder.Username,
				DisplayName: highestBid.Bidder.DisplayName,
			},
			Item: &itemResponse{
				Name:        highestBid.Item.Name,
				Description: highestBid.Item.Description,
				ImageRef:    highestBid.Item.ImageRef,
			},
		}
	}

	rawResponse, err := json.Marshal(responseObjects)
	if err != nil {
		return errors.Wrap(err, "unable to marshal response")
	}
	fmt.Fprint(w, string(rawResponse))
	return nil
}

// PostBid is the handler that lets a user bid for an existing item.
func (handler *AuctionHandler) PostBid(w http.ResponseWriter, r *http.Request) error {
	itemName := mux.Vars(r)["itemName"]
	username := auth.ExtractUsername(r.Context())

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "could not read body")
	}

	var request postBidRequest
	err = json.Unmarshal(rawBody, &request)
	if err != nil {
		log.Info(r.Context(), "invalid json request", "body", string(rawBody), "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	defer r.Body.Close()

	user, err := handler.userClient.Get(r.Context(), username)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
			response, marshalErr := json.Marshal(newErrorResponse("user does not exist"))
			if marshalErr != nil {
				return errors.Wrap(marshalErr, "could not marshal error response")
			}
			fmt.Fprint(w, string(response))
			return nil
		}
		return errors.Wrap(err, "could not retrieve user")
	}

	item, err := handler.auctionItemClient.Get(r.Context(), itemName)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
			response, marshalErr := json.Marshal(newErrorResponse("item does not exist"))
			if marshalErr != nil {
				return errors.Wrap(marshalErr, "could not marshal error response")
			}
			fmt.Fprint(w, string(response))
			return nil
		}
		return errors.Wrap(err, "could not retrieve item")
	}

	_, err = handler.auctionBidClient.PlaceBid(r.Context(), user, item, request.BidAmount)
	if err != nil {
		if errors.Is(err, storage.ErrBidTooLow) {
			w.WriteHeader(http.StatusBadRequest)
			response, marshalErr := json.Marshal(newErrorResponse("bid too low"))
			if marshalErr != nil {
				return errors.Wrap(marshalErr, "could not marshal error response")
			}
			fmt.Fprint(w, string(response))
			return nil
		}
		return errors.Wrap(err, "could not place bid")
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}
