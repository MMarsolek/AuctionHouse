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

// ----- Start Documentation Generation Types --------------

// getItemRequestDoc is for swagger generation only.
// swagger:parameters getItemRequest
type getItemRequestDoc struct {
	// Name of the item.
	//
	// In: path
	Name string `json:"itemName"`

	// Expected to be "Bearer <auth_token>"
	//
	// In: header
	Authentication string
}

// Contains data about the item and how to identify them.
//
// swagger:response getItemResponse
type getItemResponseDoc struct {

	// In: body
	Body struct {
		// The name of the item.
		//
		// Required: true
		Name string `json:"name"`

		// The reference to the image source.
		ImageRef string `json:"image,omitempty"`

		// The description of the item.
		Description string `json:"description,omitempty"`
	}
}

// ----- End Documentation Generation Types --------------

// GetItem is the handler that retrieves the model.AuctionItem as serialized JSON.
//
// swagger:route GET /api/v1/auctions/items/{itemName} Auctions getItemRequest
//
// Gets the item specified by the name.
//
// This will retrieve a item from storage based on the name.
//
//  Produces:
//  - application/json
//
//  Schemes: http
//
//  Security:
//    api_key:
//
//  Responses:
//    200: getItemResponse
//    404: noBody
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

// ----- Start Documentation Generation Types --------------

// getItemsRequestDoc is for swagger generation only.
// swagger:parameters getItemsRequest
type getItemsRequestDoc struct {
	// Expected to be "Bearer <auth_token>"
	//
	// In: header
	Authentication string
}

// Contains data about the item and how to identify them.
//
// swagger:response getItemsResponse
type getItemsResponseDoc struct {

	// In: body
	Body []struct {
		// The name of the item.
		//
		// Required: true
		Name string `json:"name"`

		// The reference to the image source.
		ImageRef string `json:"image,omitempty"`

		// The description of the item.
		Description string `json:"description,omitempty"`
	}
}

// ----- End Documentation Generation Types --------------

// GetItems is the handler that retrieves all model.AuctionItems as serialized JSON.
//
// swagger:route GET /api/v1/auctions/items Auctions getItemsRequest
//
// Gets all items that are currently stored in the system.
//
// This will retrieve all items from storage.
//
//  Produces:
//  - application/json
//
//  Schemes: http
//
//  Security:
//    api_key:
//
//  Responses:
//    200: getItemsResponse
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

// ----- Start Documentation Generation Types --------------

// postItemRequestDoc is for swagger generation only.
// swagger:parameters postItemRequest
type postItemRequestDoc struct {
	// Expected to be "Bearer <auth_token>"
	//
	// In: header
	Authentication string

	// In: body
	Body struct {
		// Name used to identify the item later.
		//
		// Required: true
		Name string `json:"name"`

		// Description of the item.
		Description string `json:"description,omitempty"`

		// Reference to the image source.
		ImageRef string `json:"image,omitempty"`
	}
}

// ----- End Documentation Generation Types --------------

// PostItem is the handler that creates a new model.AuctionItem.
//
// swagger:route POST /api/v1/auctions/items Auctions postItemRequest
//
// Creates a new item for auction.
//
// This will create a new item available for being auctioned. This route is only available to Admin users.
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
//
//  Schemes: http
//
//  Security:
//    api_key:
//
//  Responses:
//    201: noBody
//    400: errorMessage
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

// ----- Start Documentation Generation Types --------------

// putItemRequestDoc is for swagger generation only.
// swagger:parameters putItemRequest
type putItemRequestDoc struct {
	// Expected to be "Bearer <auth_token>"
	//
	// In: header
	Authentication string

	// In: path
	ItemName string `json:"itemName"`

	// In: body
	Body struct {
		// Name used to identify the item later.
		//
		// Required: true
		Name string `json:"name"`

		// Description of the item.
		Description string `json:"description,omitempty"`

		// Reference to the image source.
		ImageRef string `json:"image,omitempty"`
	}
}

// ----- End Documentation Generation Types --------------

// PutItem is the handler that updates the fields of a model.AuctionItem.
//
// swagger:route PUT /api/v1/auctions/items/{itemName} Auctions putItemRequest
//
// Updates fields for an item.
//
// This will update an existing item available for being auctioned. This route is only available to Admin users.
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
//
//  Schemes: http
//
//  Security:
//    api_key:
//
//  Responses:
//    200: noBody
//    400: errorMessage
//    404: errorMessage
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

// ----- Start Documentation Generation Types --------------

// deleteItemRequestDoc is for swagger generation only.
// swagger:parameters deleteItemRequest
type deleteItemRequestDoc struct {
	// Expected to be "Bearer <auth_token>"
	//
	// In: header
	Authentication string

	// In: path
	ItemName string `json:"itemName"`
}

// ----- End Documentation Generation Types --------------

// DeleteItem is the handler that removes a model.AuctionItem from storage.
//
// swagger:route DELETE /api/v1/auctions/items/{itemName} Auctions deleteItemRequest
//
// Deletes an item from the server.
//
// This will delete an existing item. This route is only available to Admin users.
//
//  Produces:
//  - application/json
//
//  Schemes: http
//
//  Security:
//    api_key:
//
//  Responses:
//    200: noBody
//    404: errorMessage
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

// ----- Start Documentation Generation Types --------------

// getHighestBidRequestDoc is for swagger generation only.
// swagger:parameters getHighestBidRequest
type getHighestBidRequestDoc struct {
	// Expected to be "Bearer <auth_token>"
	//
	// In: header
	Authentication string

	// In: path
	ItemName string `json:"itemName"`
}

// Contains data about the bid, what the item is, and who made it.
//
// swagger:response getHighestBidResponse
type getHighestBidResponseDoc struct {

	// In: body
	Body struct {
		// The amount of money being bid for this item.
		//
		// Required: true
		BidAmount int `json:"bidAmount"`

		// The item being bid on.
		//
		// Required: true
		Item struct {
			// The name of the item.
			//
			// Required: true
			Name string `json:"name"`

			// The reference to the image source.
			ImageRef string `json:"image,omitempty"`

			// The description of the item.
			Description string `json:"description,omitempty"`
		} `json:"item"`

		// The user who bid on the item.
		//
		// Required: true
		Bidder struct {
			// The username of the user.
			//
			// Required: true
			Username string `json:"username"`

			// The display name of the user.
			DisplayName string `json:"displayName,omitempty"`
		} `json:"bidder"`
	}
}

// ----- End Documentation Generation Types --------------

// GetHighestBid is the handler that finds the highest bid for a model.AuctionItem and retrieves the model.AuctionBid as
// serialized JSON.
//
// swagger:route GET /api/v1/auctions/bids/{itemName} Auctions getHighestBidRequest
//
// Retrieves the highest bid for the specified item.
//
// This will retrieve the highest bid for the specified item.
//
//  Produces:
//  - application/json
//
//  Schemes: http
//
//  Security:
//    api_key:
//
//  Responses:
//    200: getHighestBidResponse
//    404: errorMessage
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

// ----- Start Documentation Generation Types --------------

// getHighestBidsRequestDoc is for swagger generation only.
// swagger:parameters getHighestBidsRequest
type getHighestBidsRequestDoc struct {
	// Expected to be "Bearer <auth_token>"
	//
	// In: header
	Authentication string
}

// Contains data about the bid, what the item is, and who made it.
//
// swagger:response getHighestBidsResponse
type getHighestBidsResponseDoc struct {

	// In: body
	Body []struct {
		// The amount of money being bid for this item.
		//
		// Required: true
		BidAmount int `json:"bidAmount"`

		// The item being bid on.
		//
		// Required: true
		Item struct {
			// The name of the item.
			//
			// Required: true
			Name string `json:"name"`

			// The reference to the image source.
			ImageRef string `json:"image,omitempty"`

			// The description of the item.
			Description string `json:"description,omitempty"`
		} `json:"item"`

		// The user who bid on the item.
		//
		// Required: true
		Bidder struct {
			// The username of the user.
			//
			// Required: true
			Username string `json:"username"`

			// The display name of the user.
			DisplayName string `json:"displayName,omitempty"`
		} `json:"bidder"`
	}
}

// ----- End Documentation Generation Types --------------

// GetHighestBid is the handler that finds the highest bids all model.AuctionItems and retrieves all model.AuctionBids as
// serialized JSON.
//
// swagger:route GET /api/v1/auctions/bids Auctions getHighestBidsRequest
//
// Retrieves the highest bid for all items.
//
// This will retrieve the highest bid for all items.
//
//  Produces:
//  - application/json
//
//  Schemes: http
//
//  Security:
//    api_key:
//
//  Responses:
//    200: getHighestBidResponse
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

// ----- Start Documentation Generation Types --------------

// postBidRequestDoc is for swagger generation only.
// swagger:parameters postBidRequest
type postBidRequestDoc struct {
	// Expected to be "Bearer <auth_token>"
	//
	// In: header
	Authentication string

	// In: path
	ItemName string `json:"itemName"`

	// In: body
	Body struct {
		// The amount to bid on the item for.
		//
		// Required: true
		BidAmount int `json:"bidAmount"`
	}
}

// ----- End Documentation Generation Types --------------

// PostBid is the handler that lets a user bid for an existing item.
//
// swagger:route POST /api/v1/auctions/bids/{itemName} Auctions postBidRequest
//
// Makes a new bid on an item.
//
// This will place a bid on the specified item. The user is identified by the authentication token. This is only
// available for Bidder users.
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
//
//  Schemes: http
//
//  Security:
//    api_key:
//
//  Responses:
//    201: noBody
//    400: errorMessage
//    404: errorMessage
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
