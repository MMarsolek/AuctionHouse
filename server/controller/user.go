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
	postUserRequest struct {
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
		Password    string `json:"password"`
	}

	postLoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	postLoginResponse struct {
		Username    string                `json:"username"`
		DisplayName string                `json:"displayName"`
		Permission  model.PermissionLevel `json:"permission"`
		AuthToken   string                `json:"authToken"`
	}
)

// UserHandler provides handlers for endpoints involving model.Users.
type UserHandler struct {
	userClient storage.UserClient
}

// NewUserHandler creates a new UserHandler with the necessary storage objects.
func NewUserHandler(userClient storage.UserClient) *UserHandler {
	return &UserHandler{
		userClient: userClient,
	}
}

// RegisterRoutes registers all of the paths to the handler functions.
func (handler *UserHandler) RegisterRoutes(router *mux.Router) {
	usersRouter := router.PathPrefix("/v1/users").Subrouter()
	usersRouterWithAuth := usersRouter.NewRoute().Subrouter()
	usersRouterWithAuth.Use(middleware.VerifyAuthToken)

	usersRouter.HandleFunc("/login", wrapHandler(handler.PostLogin)).Methods(http.MethodPost)
	usersRouter.HandleFunc("", wrapHandler(handler.PostUser)).Methods(http.MethodPost)
	usersRouterWithAuth.HandleFunc("/{username}", wrapHandler(handler.GetUser)).Methods(http.MethodGet)
}

// GetUser is the handler that retrieves the model.User as serialized JSON.
func (handler *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) error {
	username := mux.Vars(r)["username"]
	r = r.WithContext(log.WithFields(r.Context(), "username", username))
	user, err := handler.userClient.Get(r.Context(), username)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return nil
		}
		return errors.Wrap(err, "could not retrieve user")
	}

	rawUser, err := json.Marshal(struct {
		Username    string                `json:"username"`
		DisplayName string                `json:"displayName"`
		Permission  model.PermissionLevel `json:"permission"`
	}{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Permission:  user.Permission,
	})

	if err != nil {
		return errors.Wrap(err, "could not marshal user")
	}

	fmt.Fprint(w, string(rawUser))
	return nil
}

// PostUser is the handler that adds a new model.User to storage.
func (handler *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) error {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "could not read body")
	}

	var request postUserRequest
	err = json.Unmarshal(rawBody, &request)
	if err != nil {
		log.Info(r.Context(), "invalid json request", "body", string(rawBody), "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	defer r.Body.Close()

	hashedPassword, err := auth.GenerateEncodedPassword(request.Password)
	if err != nil {
		if errors.Is(err, auth.ErrEmptyPassword) {
			w.WriteHeader(http.StatusBadRequest)
			response, marshalErr := json.Marshal(newErrorResponse("invalid password"))
			if marshalErr != nil {
				return errors.Wrap(marshalErr, "could not marshal error response")
			}
			fmt.Fprint(w, string(response))
			return nil
		}
		return errors.Wrap(err, "could not generate password")
	}

	newUser := &model.User{
		Username:       request.Username,
		DisplayName:    request.DisplayName,
		HashedPassword: hashedPassword,
		Permission:     model.PermissionLevelBidder,
	}

	err = handler.userClient.Create(r.Context(), newUser)
	if err != nil {
		if errors.Is(err, storage.ErrEntityAlreadyExists) {
			w.WriteHeader(http.StatusBadRequest)
			response, marshalErr := json.Marshal(newErrorResponse("user already exists"))
			if marshalErr != nil {
				return errors.Wrap(marshalErr, "could not marshal error response")
			}
			fmt.Fprint(w, string(response))
			return nil
		}
		return errors.Wrap(err, "could not store user")
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

// PostLogin is the handler that retrieves an auth token for the specified model.User.
func (handler *UserHandler) PostLogin(w http.ResponseWriter, r *http.Request) error {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "could not read body")
	}

	var request postLoginRequest
	err = json.Unmarshal(rawBody, &request)
	if err != nil {
		log.Info(r.Context(), "invalid json request", "body", string(rawBody), "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	defer r.Body.Close()

	user, err := handler.userClient.Get(r.Context(), request.Username)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
			response, marshalErr := json.Marshal(newErrorResponse("password mismatch or the user does not exist"))
			if marshalErr != nil {
				return errors.Wrap(marshalErr, "could not marshal error response")
			}
			fmt.Fprint(w, string(response))
			return nil
		}
		return errors.Wrap(err, "unable to retrieve user")
	}

	match, err := auth.ComparePasswordAndHash(request.Password, user.HashedPassword)
	if err != nil {
		return errors.Wrap(err, "unable to verify password hash")
	}

	if !match {
		w.WriteHeader(http.StatusBadRequest)
		response, marshalErr := json.Marshal(newErrorResponse("password mismatch or the user does not exist"))
		if marshalErr != nil {
			return errors.Wrap(marshalErr, "could not marshal error response")
		}
		fmt.Fprint(w, string(response))
		return nil
	}

	token, err := auth.NewToken(user)
	if err != nil {
		return errors.Wrapf(err, "unable to generate new token for %s", user.Username)
	}

	rawResponse, err := json.Marshal(postLoginResponse{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Permission:  user.Permission,
		AuthToken:   string(token),
	})
	if err != nil {
		return errors.Wrap(err, "unable to marshal login response")
	}

	fmt.Fprint(w, string(rawResponse))
	return nil
}
