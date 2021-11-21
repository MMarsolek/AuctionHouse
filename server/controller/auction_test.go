package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/MMarsolek/AuctionHouse/storage"
	"github.com/MMarsolek/AuctionHouse/storage/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type auctionHandlerTestSuite struct {
	suite.Suite

	client          *http.Client
	server          *httptest.Server
	userStoreMock   *mocks.UserClient
	auctionItemMock *mocks.AuctionItemClient
	auctionBidMock  *mocks.AuctionBidClient
	handler         *AuctionHandler
}

func (ts *auctionHandlerTestSuite) SetupSuite() {
	ts.handler = NewAuctionHandler(ts.userStoreMock, ts.auctionItemMock, ts.auctionBidMock)
	var router *mux.Router
	ts.server, router = newTestServer()
	ts.handler.RegisterRoutes(router)
	ts.client = &http.Client{}
}

func (ts *auctionHandlerTestSuite) SetupTest() {
	ts.userStoreMock = new(mocks.UserClient)
	ts.auctionItemMock = new(mocks.AuctionItemClient)
	ts.auctionBidMock = new(mocks.AuctionBidClient)
	ts.handler.userClient = ts.userStoreMock
	ts.handler.auctionItemClient = ts.auctionItemMock
	ts.handler.auctionBidClient = ts.auctionBidMock
}

func (ts *auctionHandlerTestSuite) TearDownTest() {
	ts.userStoreMock.AssertExpectations(ts.T())
	ts.auctionItemMock.AssertExpectations(ts.T())
	ts.auctionBidMock.AssertExpectations(ts.T())
}

func (ts *auctionHandlerTestSuite) TearDownSuite() {
	ts.server.Close()
}

func TestAuctionHandler(t *testing.T) {
	suite.Run(t, new(auctionHandlerTestSuite))
}

func (ts *auctionHandlerTestSuite) TestGetItemRetrievesItem() {
	getItemTest := func(permission model.PermissionLevel) {
		item := &model.AuctionItem{
			Name:        "foo",
			ImageRef:    "ref",
			Description: "description",
		}
		ts.auctionItemMock.On("Get", mock.AnythingOfType("*context.valueCtx"), item.Name).Return(item, nil)

		r := ts.makeAuthenticatedRequest(http.MethodGet, fmt.Sprintf("items/%s", item.Name), nil, &model.User{
			Permission: model.PermissionLevelBidder,
		})
		response, err := ts.client.Do(r)
		ts.Require().NoError(err)

		defer response.Body.Close()
		ts.Require().EqualValues(http.StatusOK, response.StatusCode)
		rawResponse, err := io.ReadAll(response.Body)
		ts.Require().NoError(err)
		var returnedItem getItemResponse
		ts.Require().NoError(json.Unmarshal(rawResponse, &returnedItem))

		ts.Require().EqualValues(item.Name, returnedItem.Name)
		ts.Require().EqualValues(item.ImageRef, returnedItem.ImageRef)
		ts.Require().EqualValues(item.Description, returnedItem.Description)
	}

	getItemTest(model.PermissionLevelAdmin)
	getItemTest(model.PermissionLevelBidder)
}

func (ts *auctionHandlerTestSuite) TestGetItem404OnNonExistantItem() {
	ts.auctionItemMock.On("Get", mock.AnythingOfType("*context.valueCtx"), "whatever").Return(nil, storage.ErrEntityNotFound)

	r := ts.makeAuthenticatedRequest(http.MethodGet, "items/whatever", nil, &model.User{
		Permission: model.PermissionLevelBidder,
	})
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)

	defer response.Body.Close()
	ts.Require().EqualValues(http.StatusNotFound, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestGetItemsRetrievesMultipleItems() {
	getItemTest := func(permission model.PermissionLevel) {
		items := []*model.AuctionItem{
			{
				Name:        "foo",
				ImageRef:    "ref",
				Description: "description",
			},
			{
				Name:        "bar",
				ImageRef:    "barref",
				Description: "bardescription",
			},
		}
		ts.auctionItemMock.On("GetAll", mock.AnythingOfType("*context.valueCtx")).Return(items, nil)

		r := ts.makeAuthenticatedRequest(http.MethodGet, "items", nil, &model.User{
			Permission: model.PermissionLevelBidder,
		})
		response, err := ts.client.Do(r)
		ts.Require().NoError(err)

		defer response.Body.Close()
		ts.Require().EqualValues(http.StatusOK, response.StatusCode)
		rawResponse, err := io.ReadAll(response.Body)
		ts.Require().NoError(err)
		var returnedItems []*getItemResponse
		ts.Require().NoError(json.Unmarshal(rawResponse, &returnedItems))

		for i, item := range items {
			ts.Require().EqualValues(item.Name, returnedItems[i].Name)
			ts.Require().EqualValues(item.ImageRef, returnedItems[i].ImageRef)
			ts.Require().EqualValues(item.Description, returnedItems[i].Description)
		}
	}

	getItemTest(model.PermissionLevelAdmin)
	getItemTest(model.PermissionLevelBidder)
}

func (ts *auctionHandlerTestSuite) TestPostItemStoresNewItem() {
	ts.auctionItemMock.On("Create", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*model.AuctionItem")).Return(nil)
	itemRequest := postItemRequest{
		Name:        "item",
		ImageRef:    "image",
		Description: "description",
	}
	rawRequest, err := json.Marshal(itemRequest)
	ts.Require().NoError(err)

	r := ts.makeAuthenticatedRequest(http.MethodPost, "items", bytes.NewReader(rawRequest), &model.User{
		Permission: model.PermissionLevelAdmin,
	})
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusCreated, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestPostItem403OnBidderRequest() {
	r := ts.makeAuthenticatedRequest(http.MethodPost, "items", nil, &model.User{
		Permission: model.PermissionLevelBidder,
	})
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusForbidden, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestPostItem400OnItemThatAlreadyExists() {
	ts.auctionItemMock.On("Create", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*model.AuctionItem")).Return(storage.ErrEntityAlreadyExists)
	itemRequest := postItemRequest{
		Name:        "item",
		ImageRef:    "image",
		Description: "description",
	}
	rawRequest, err := json.Marshal(itemRequest)
	ts.Require().NoError(err)

	r := ts.makeAuthenticatedRequest(http.MethodPost, "items", bytes.NewReader(rawRequest), &model.User{
		Permission: model.PermissionLevelAdmin,
	})
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusBadRequest, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestPutItemStoresUpdates() {
	ts.auctionItemMock.On("Update", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*model.AuctionItem")).Return(nil)
	itemRequest := putItemRequest{
		ImageRef:    "image",
		Description: "description",
	}
	rawRequest, err := json.Marshal(itemRequest)
	ts.Require().NoError(err)

	r := ts.makeAuthenticatedRequest(http.MethodPut, "items/someItem", bytes.NewReader(rawRequest), &model.User{
		Permission: model.PermissionLevelAdmin,
	})
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusOK, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestPutItem403OnBidderRequest() {
	r := ts.makeAuthenticatedRequest(http.MethodPut, "items/someItem", nil, &model.User{
		Permission: model.PermissionLevelBidder,
	})
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusForbidden, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestPutItem404WhenItemNotFound() {
	ts.auctionItemMock.On("Update", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*model.AuctionItem")).Return(storage.ErrEntityNotFound)
	itemRequest := putItemRequest{
		ImageRef:    "image",
		Description: "description",
	}
	rawRequest, err := json.Marshal(itemRequest)
	ts.Require().NoError(err)

	r := ts.makeAuthenticatedRequest(http.MethodPut, "items/someItem", bytes.NewReader(rawRequest), &model.User{
		Permission: model.PermissionLevelAdmin,
	})
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusNotFound, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestDeleteItemRemovesItemFromStorage() {
	itemName := "someItem"
	ts.auctionItemMock.On("Delete", mock.AnythingOfType("*context.valueCtx"), itemName).Return(nil)

	r := ts.makeAuthenticatedRequest(http.MethodDelete, fmt.Sprintf("items/%s", itemName), nil, &model.User{
		Permission: model.PermissionLevelAdmin,
	})
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusOK, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestDeleteItem403OnBidderRequest() {
	r := ts.makeAuthenticatedRequest(http.MethodDelete, "items/someItem", nil, &model.User{
		Permission: model.PermissionLevelBidder,
	})
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusForbidden, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestDeleteItem404WhenItemDoesNotExist() {
	itemName := "someItem"
	ts.auctionItemMock.On("Delete", mock.AnythingOfType("*context.valueCtx"), itemName).Return(storage.ErrEntityNotFound)

	r := ts.makeAuthenticatedRequest(http.MethodDelete, fmt.Sprintf("items/%s", itemName), nil, &model.User{
		Permission: model.PermissionLevelAdmin,
	})
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusNotFound, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestGetHighestBidRetrievesBidFromStorage() {
	getHighestBidTest := func(permission model.PermissionLevel) {
		item := &model.AuctionItem{
			Name: "item",
		}
		ts.auctionItemMock.On("Get", mock.AnythingOfType("*context.valueCtx"), item.Name).Return(item, nil)
		highestBid := &model.AuctionBid{
			BidAmount: 100,
			Item:      item,
			Bidder: &model.User{
				Username: "user",
			},
		}
		ts.auctionBidMock.On("GetHighestBid", mock.AnythingOfType("*context.valueCtx"), item).Return(highestBid, nil)

		r := ts.makeAuthenticatedRequest(http.MethodGet, fmt.Sprintf("bids/%s", item.Name), nil, &model.User{
			Permission: permission,
		})
		response, err := ts.client.Do(r)
		ts.Require().NoError(err)
		ts.Require().EqualValues(http.StatusOK, response.StatusCode)

		defer response.Body.Close()
		rawResponse, err := io.ReadAll(response.Body)
		ts.Require().NoError(err)
		var bidResponse getHighestBidsResponse
		ts.Require().NoError(json.Unmarshal(rawResponse, &bidResponse))

		ts.Require().EqualValues(highestBid.BidAmount, bidResponse.BidAmount)
		ts.Require().EqualValues(highestBid.Item, bidResponse.Item)
		ts.Require().EqualValues(highestBid.Bidder.Username, bidResponse.Bidder.Username)
	}

	getHighestBidTest(model.PermissionLevelAdmin)
	getHighestBidTest(model.PermissionLevelBidder)
}

func (ts *auctionHandlerTestSuite) TestGetHighestBid404WhenItemNotFound() {
	getHighestBidTest := func(permission model.PermissionLevel) {
		itemName := "someItem"
		ts.auctionItemMock.On("Get", mock.AnythingOfType("*context.valueCtx"), itemName).Return(nil, storage.ErrEntityNotFound)
		r := ts.makeAuthenticatedRequest(http.MethodGet, fmt.Sprintf("bids/%s", itemName), nil, &model.User{
			Permission: permission,
		})
		response, err := ts.client.Do(r)
		ts.Require().NoError(err)
		ts.Require().EqualValues(http.StatusNotFound, response.StatusCode)
	}

	getHighestBidTest(model.PermissionLevelAdmin)
	getHighestBidTest(model.PermissionLevelBidder)
}

func (ts *auctionHandlerTestSuite) TestGetHighestBid404WhenItemHasNoBids() {
	getHighestBidTest := func(permission model.PermissionLevel) {
		item := &model.AuctionItem{
			Name: "item",
		}
		ts.auctionItemMock.On("Get", mock.AnythingOfType("*context.valueCtx"), item.Name).Return(item, nil)
		ts.auctionBidMock.On("GetHighestBid", mock.AnythingOfType("*context.valueCtx"), item).Return(nil, storage.ErrEntityNotFound)

		r := ts.makeAuthenticatedRequest(http.MethodGet, fmt.Sprintf("bids/%s", item.Name), nil, &model.User{
			Permission: permission,
		})
		response, err := ts.client.Do(r)
		ts.Require().NoError(err)
		ts.Require().EqualValues(http.StatusNotFound, response.StatusCode)
	}

	getHighestBidTest(model.PermissionLevelAdmin)
	getHighestBidTest(model.PermissionLevelBidder)
}

func (ts *auctionHandlerTestSuite) TestGetHighestBidsRetrievesAllBidsForEveryItem() {
	getHighestBidsTest := func(permission model.PermissionLevel) {
		items := []*model.AuctionBid{
			{
				BidAmount: 100,
				Item: &model.AuctionItem{
					Name: "item1",
				},
				Bidder: &model.User{
					Username: "user1",
				},
			},
			{
				BidAmount: 150,
				Item: &model.AuctionItem{
					Name: "item2",
				},
				Bidder: &model.User{
					Username: "user1",
				},
			},
		}
		ts.auctionBidMock.On("GetAllHighestBids", mock.AnythingOfType("*context.valueCtx")).Return(items, nil)

		r := ts.makeAuthenticatedRequest(http.MethodGet, "bids", nil, &model.User{
			Permission: permission,
		})
		response, err := ts.client.Do(r)
		ts.Require().NoError(err)
		ts.Require().EqualValues(http.StatusOK, response.StatusCode)

		defer response.Body.Close()
		rawResponse, err := io.ReadAll(response.Body)
		ts.Require().NoError(err)
		var highestBids []*getHighestBidsResponse
		ts.Require().NoError(json.Unmarshal(rawResponse, &highestBids))

		for i, item := range items {
			ts.Require().EqualValues(item.BidAmount, highestBids[i].BidAmount)
			ts.Require().EqualValues(item.Item.Name, highestBids[i].Item.Name)
			ts.Require().EqualValues(item.Bidder.Username, highestBids[i].Bidder.Username)
		}
	}

	getHighestBidsTest(model.PermissionLevelAdmin)
	getHighestBidsTest(model.PermissionLevelBidder)
}

func (ts *auctionHandlerTestSuite) TestPostBidPlacesBidOnItem() {
	user := &model.User{
		Username:    "user1",
		DisplayName: "User 1",
		Permission:  model.PermissionLevelBidder,
	}
	ts.userStoreMock.On("Get", mock.AnythingOfType("*context.valueCtx"), user.Username).Return(user, nil)
	item := &model.AuctionItem{
		Name:        "Item1",
		ImageRef:    "image",
		Description: "desc",
	}
	bidAmount := 100
	ts.auctionItemMock.On("Get", mock.AnythingOfType("*context.valueCtx"), item.Name).Return(item, nil)
	ts.auctionBidMock.On("PlaceBid", mock.AnythingOfType("*context.valueCtx"), user, item, bidAmount).Return(nil, nil)

	rawRequest, err := json.Marshal(postBidRequest{
		BidAmount: bidAmount,
	})
	ts.Require().NoError(err)
	r := ts.makeAuthenticatedRequest(http.MethodPost, fmt.Sprintf("bids/%s", item.Name), bytes.NewReader(rawRequest), user)
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusCreated, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestPostBid404WhenUserNotFound() {
	user := &model.User{
		Username:    "user1",
		DisplayName: "User 1",
		Permission:  model.PermissionLevelBidder,
	}
	ts.userStoreMock.On("Get", mock.AnythingOfType("*context.valueCtx"), user.Username).Return(nil, storage.ErrEntityNotFound)

	rawRequest, err := json.Marshal(postBidRequest{
		BidAmount: 100,
	})
	ts.Require().NoError(err)
	r := ts.makeAuthenticatedRequest(http.MethodPost, "bids/doesnotmatter", bytes.NewReader(rawRequest), user)
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusNotFound, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestPostBid404WhenItemNotFound() {
	user := &model.User{
		Username:    "user1",
		DisplayName: "User 1",
		Permission:  model.PermissionLevelBidder,
	}
	ts.userStoreMock.On("Get", mock.AnythingOfType("*context.valueCtx"), user.Username).Return(user, nil)
	item := &model.AuctionItem{
		Name:        "Item1",
		ImageRef:    "image",
		Description: "desc",
	}
	ts.auctionItemMock.On("Get", mock.AnythingOfType("*context.valueCtx"), item.Name).Return(nil, storage.ErrEntityNotFound)

	rawRequest, err := json.Marshal(postBidRequest{
		BidAmount: 100,
	})
	ts.Require().NoError(err)
	r := ts.makeAuthenticatedRequest(http.MethodPost, fmt.Sprintf("bids/%s", item.Name), bytes.NewReader(rawRequest), user)
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusNotFound, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestPostBid400WhenBidTooLow() {
	user := &model.User{
		Username:    "user1",
		DisplayName: "User 1",
		Permission:  model.PermissionLevelBidder,
	}
	ts.userStoreMock.On("Get", mock.AnythingOfType("*context.valueCtx"), user.Username).Return(user, nil)
	item := &model.AuctionItem{
		Name:        "Item1",
		ImageRef:    "image",
		Description: "desc",
	}
	bidAmount := 100
	ts.auctionItemMock.On("Get", mock.AnythingOfType("*context.valueCtx"), item.Name).Return(item, nil)
	ts.auctionBidMock.On("PlaceBid", mock.AnythingOfType("*context.valueCtx"), user, item, bidAmount).Return(nil, storage.ErrBidTooLow)

	rawRequest, err := json.Marshal(postBidRequest{
		BidAmount: bidAmount,
	})
	ts.Require().NoError(err)
	r := ts.makeAuthenticatedRequest(http.MethodPost, fmt.Sprintf("bids/%s", item.Name), bytes.NewReader(rawRequest), user)
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusBadRequest, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) TestPostBid403OnAdminRequest() {
	r := ts.makeAuthenticatedRequest(http.MethodPost, "bids/someItem", nil, &model.User{
		Permission: model.PermissionLevelAdmin,
	})
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusForbidden, response.StatusCode)
}

func (ts *auctionHandlerTestSuite) fullPath(path string) string {
	return fmt.Sprintf("%s/api/v1/auctions/%s", ts.server.URL, path)
}

func (ts *auctionHandlerTestSuite) makeAuthenticatedRequest(method string, path string, body io.Reader, user *model.User) *http.Request {
	return makeAuthenticatedRequest(ts.T(), method, ts.fullPath(path), body, user)
}
