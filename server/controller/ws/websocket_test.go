package ws

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/MMarsolek/AuctionHouse/server/controller/auth"
	"github.com/MMarsolek/AuctionHouse/server/controller/middleware"
	"github.com/MMarsolek/AuctionHouse/storage"
	"github.com/MMarsolek/AuctionHouse/storage/mocks"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	testUserName = "bidder"
	testItemName = "item1"
)

type handlerTestSuite struct {
	suite.Suite

	server   *httptest.Server
	handler  *Handler
	userMock *mocks.UserClient
	itemMock *mocks.AuctionItemClient
	bidMock  *mocks.AuctionBidClient
}

func (ts *handlerTestSuite) SetupSuite() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	ts.server = httptest.NewServer(middleware.RemoveTrailingSlash(router))
	ts.handler = NewHandler(ts.userMock, ts.itemMock, ts.bidMock)
	ts.handler.RegisterRoutes(router)
}

func (ts *handlerTestSuite) SetupTest() {
	ts.userMock = new(mocks.UserClient)
	ts.itemMock = new(mocks.AuctionItemClient)
	ts.bidMock = new(mocks.AuctionBidClient)
	ts.handler.userClient = ts.userMock
	ts.handler.itemClient = ts.itemMock
	ts.handler.bidClient = ts.bidMock
}

func (ts *handlerTestSuite) TearDownTest() {
	ts.userMock.AssertExpectations(ts.T())
	ts.itemMock.AssertExpectations(ts.T())
	ts.bidMock.AssertExpectations(ts.T())
}

func (ts *handlerTestSuite) TearDownSuite() {
	ts.server.Close()
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}

func (ts *handlerTestSuite) TestServeWSHandlesMultipleConnections() {
	for i := 0; i < 10; i++ {
		defer ts.createWebsocket().Close()
	}
}

func (ts *handlerTestSuite) TestServeWSCanPlaceBid() {
	user := &model.User{
		Username: testUserName,
	}
	item := &model.AuctionItem{
		Name: testItemName,
	}
	bidAmount := 1000
	ts.userMock.On("Get", mock.AnythingOfType("*context.valueCtx"), testUserName).Return(user, nil)
	ts.itemMock.On("Get", mock.AnythingOfType("*context.valueCtx"), testItemName).Return(item, nil)
	ts.bidMock.On("PlaceBid", mock.AnythingOfType("*context.valueCtx"), user, item, bidAmount).Return(nil, nil)
	ws := ts.createWebsocket()
	defer ws.Close()

	ts.Require().NoError(ws.WriteJSON(commandMessage{
		Command: SocketCommandPlaceBid,
		Payload: commandMessagePlaceBid{
			ItemName:  item.Name,
			BidAmount: bidAmount,
		},
	}))
	var response responseMessage
	ts.Require().NoError(ws.ReadJSON(&response))
	ts.Require().EqualValues(SocketCommandPlaceBid, response.Command)
	ts.Require().EqualValues(http.StatusCreated, response.StatusCode)
	ts.Require().NotNil(response.Data)
	ts.Require().IsType((map[string]interface{})(nil), response.Data)
	result := response.Data.(map[string]interface{})
	ts.Require().EqualValues(result["itemName"], item.Name)
	ts.Require().EqualValues(result["username"], user.Username)
	ts.Require().EqualValues(result["amount"], bidAmount)
}

func (ts *handlerTestSuite) TestServeWSReturnsErrorJSONOnItemNotFound() {
	user := &model.User{
		Username: testUserName,
	}
	item := &model.AuctionItem{
		Name: testItemName,
	}
	bidAmount := 1000
	ts.userMock.On("Get", mock.AnythingOfType("*context.valueCtx"), testUserName).Return(user, nil)
	ts.itemMock.On("Get", mock.AnythingOfType("*context.valueCtx"), testItemName).Return(item, nil)
	ts.bidMock.On("PlaceBid", mock.AnythingOfType("*context.valueCtx"), user, item, bidAmount).Return(nil, storage.ErrBidTooLow)
	ws := ts.createWebsocket()
	defer ws.Close()

	ts.Require().NoError(ws.WriteJSON(commandMessage{
		Command: SocketCommandPlaceBid,
		Payload: commandMessagePlaceBid{
			ItemName:  item.Name,
			BidAmount: bidAmount,
		},
	}))
	var response responseMessage
	ts.Require().NoError(ws.ReadJSON(&response))
	ts.Require().EqualValues(SocketCommandPlaceBid, response.Command)
	ts.Require().EqualValues(http.StatusBadRequest, response.StatusCode)
	ts.Require().EqualValues(fmt.Sprintf("bid amount is too low"), response.Message)
	ts.Require().Nil(response.Data)
}

func (ts *handlerTestSuite) TestServeWSReturnsErrorJSONOnBidBeingTooLow() {
	user := &model.User{
		Username: testUserName,
	}
	item := &model.AuctionItem{
		Name: testItemName,
	}
	bidAmount := 1000
	ts.userMock.On("Get", mock.AnythingOfType("*context.valueCtx"), testUserName).Return(user, nil)
	ts.itemMock.On("Get", mock.AnythingOfType("*context.valueCtx"), testItemName).Return(nil, storage.ErrEntityNotFound)
	ws := ts.createWebsocket()
	defer ws.Close()

	ts.Require().NoError(ws.WriteJSON(commandMessage{
		Command: SocketCommandPlaceBid,
		Payload: commandMessagePlaceBid{
			ItemName:  item.Name,
			BidAmount: bidAmount,
		},
	}))
	var response responseMessage
	ts.Require().NoError(ws.ReadJSON(&response))
	ts.Require().EqualValues(SocketCommandPlaceBid, response.Command)
	ts.Require().EqualValues(http.StatusNotFound, response.StatusCode)
	ts.Require().EqualValues(fmt.Sprintf("could not find item '%s'", item.Name), response.Message)
	ts.Require().Nil(response.Data)
}

func (ts *handlerTestSuite) TestServeWSFailsOnUnknownCommand() {
	ws := ts.createWebsocket()
	defer ws.Close()

	ts.Require().NoError(ws.WriteJSON(commandMessage{
		Command: SocketCommand("not a real command"),
	}))
	var response responseMessage
	ts.Require().NoError(ws.ReadJSON(&response))
	ts.Require().EqualValues(SocketCommandUnknown, response.Command)
	ts.Require().EqualValues(http.StatusBadRequest, response.StatusCode)
	ts.Require().EqualValues("invalid request format", response.Message)
}

func (ts *handlerTestSuite) TestServeWSBroadcastsBidsToAllClients() {
	user := &model.User{
		Username: testUserName,
	}
	item := &model.AuctionItem{
		Name: testItemName,
	}
	bidAmount := 1000
	ts.userMock.On("Get", mock.AnythingOfType("*context.valueCtx"), testUserName).Return(user, nil)
	ts.itemMock.On("Get", mock.AnythingOfType("*context.valueCtx"), testItemName).Return(item, nil)
	ts.bidMock.On("PlaceBid", mock.AnythingOfType("*context.valueCtx"), user, item, bidAmount).Return(nil, nil)
	sockets := make([]*websocket.Conn, 10)
	for i := range sockets {
		sockets[i] = ts.createWebsocket()
		defer sockets[i].Close()
	}

	ts.Require().NoError(sockets[0].WriteJSON(commandMessage{
		Command: SocketCommandPlaceBid,
		Payload: commandMessagePlaceBid{
			ItemName:  item.Name,
			BidAmount: bidAmount,
		},
	}))

	for _, ws := range sockets {
		var response responseMessage
		ts.Require().NoError(ws.ReadJSON(&response))
		ts.Require().EqualValues(SocketCommandPlaceBid, response.Command)
		ts.Require().EqualValues(http.StatusCreated, response.StatusCode)
		ts.Require().NotNil(response.Data)
		ts.Require().IsType((map[string]interface{})(nil), response.Data)
		result := response.Data.(map[string]interface{})
		ts.Require().EqualValues(result["itemName"], item.Name)
		ts.Require().EqualValues(result["username"], user.Username)
		ts.Require().EqualValues(result["amount"], bidAmount)
	}
}

func (ts *handlerTestSuite) TestServeWSNoLongerBroadcastsToDisconnectedClient() {
	user := &model.User{
		Username: testUserName,
	}
	item := &model.AuctionItem{
		Name: testItemName,
	}
	bidAmount := 1000
	ts.userMock.On("Get", mock.AnythingOfType("*context.valueCtx"), testUserName).Return(user, nil)
	ts.itemMock.On("Get", mock.AnythingOfType("*context.valueCtx"), testItemName).Return(item, nil)
	ts.bidMock.On("PlaceBid", mock.AnythingOfType("*context.valueCtx"), user, item, bidAmount).Return(nil, nil)
	sockets := make([]*websocket.Conn, 10)
	prematureCloseIndex := 7
	for i := range sockets {
		sockets[i] = ts.createWebsocket()
		if i != prematureCloseIndex {
			defer sockets[i].Close()
		}
	}
	ts.Require().NoError(sockets[prematureCloseIndex].Close())

	ts.Require().NoError(sockets[0].WriteJSON(commandMessage{
		Command: SocketCommandPlaceBid,
		Payload: commandMessagePlaceBid{
			ItemName:  item.Name,
			BidAmount: bidAmount,
		},
	}))

	for i, ws := range sockets {
		if i == prematureCloseIndex {
			continue
		}
		var response responseMessage
		ts.Require().NoError(ws.ReadJSON(&response))
		ts.Require().EqualValues(SocketCommandPlaceBid, response.Command)
		ts.Require().EqualValues(http.StatusCreated, response.StatusCode)
		ts.Require().NotNil(response.Data)
		ts.Require().IsType((map[string]interface{})(nil), response.Data)
		result := response.Data.(map[string]interface{})
		ts.Require().EqualValues(result["itemName"], item.Name)
		ts.Require().EqualValues(result["username"], user.Username)
		ts.Require().EqualValues(result["amount"], bidAmount)
	}

	unused := make(map[string]interface{})
	err := sockets[prematureCloseIndex].ReadJSON(&unused)
	ts.Require().Error(err)
	ts.Require().ErrorIs(err, net.ErrClosed)
}

func (ts *handlerTestSuite) createWebsocket() *websocket.Conn {
	token, err := auth.NewToken(&model.User{
		Username:   testUserName,
		Permission: model.PermissionLevelBidder,
	})
	ts.Require().NoError(err)
	ws, response, err := websocket.DefaultDialer.Dial(ts.serverURL(), http.Header(map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", token)},
	}))
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusSwitchingProtocols, response.StatusCode)

	return ws
}

func (ts *handlerTestSuite) serverURL() string {
	return fmt.Sprintf("ws://%s/api/ws", strings.TrimPrefix(ts.server.URL, "http://"))
}
