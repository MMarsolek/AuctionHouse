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
	"github.com/MMarsolek/AuctionHouse/server/controller/auth"
	"github.com/MMarsolek/AuctionHouse/storage"
	"github.com/MMarsolek/AuctionHouse/storage/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type userHandlerTestSuite struct {
	suite.Suite

	client        *http.Client
	server        *httptest.Server
	userStoreMock *mocks.UserClient
	handler       *UserHandler
}

func (ts *userHandlerTestSuite) SetupSuite() {
	ts.handler = NewUserHandler(ts.userStoreMock)
	var router *mux.Router
	ts.server, router = newTestServer()
	ts.handler.RegisterRoutes(router)
	ts.client = &http.Client{}
}

func (ts *userHandlerTestSuite) SetupTest() {
	ts.userStoreMock = new(mocks.UserClient)
	ts.handler.userClient = ts.userStoreMock
}

func (ts *userHandlerTestSuite) TearDownTest() {
	ts.userStoreMock.AssertExpectations(ts.T())
}

func (ts *userHandlerTestSuite) TearDownSuite() {
	ts.server.Close()
}

func TestUserHandler(t *testing.T) {
	suite.Run(t, new(userHandlerTestSuite))
}

func (ts *userHandlerTestSuite) TestGetUserReturnsExpectedUser() {
	user := &model.User{
		Username:    "testuser",
		DisplayName: "Test User",
		Permission:  model.PermissionLevelAdmin,
	}
	ts.userStoreMock.On("Get", mock.AnythingOfType("*context.valueCtx"), user.Username).Return(user, nil)

	r := ts.makeAuthenticatedRequest(http.MethodGet, user.Username, user)
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)

	defer response.Body.Close()
	ts.Require().EqualValues(http.StatusOK, response.StatusCode)
	rawResponse, err := io.ReadAll(response.Body)
	ts.Require().NoError(err)
	var returnedUser model.User
	ts.Require().NoError(json.Unmarshal(rawResponse, &returnedUser))

	ts.Require().EqualValues(user.Username, returnedUser.Username)
	ts.Require().EqualValues(user.DisplayName, returnedUser.DisplayName)
	ts.Require().EqualValues(user.Permission, returnedUser.Permission)
}

func (ts *userHandlerTestSuite) TestGetUser404OnUserNotFound() {
	username := "doesntexist"
	ts.userStoreMock.On("Get", mock.AnythingOfType("*context.valueCtx"), username).Return(nil, storage.ErrEntityNotFound)

	r := ts.makeAuthenticatedRequest(http.MethodGet, username, &model.User{
		Username:   username,
		Permission: model.PermissionLevelAdmin,
	})
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)

	defer response.Body.Close()
	ts.Require().EqualValues(http.StatusNotFound, response.StatusCode)
}

func (ts *userHandlerTestSuite) TestGetUser401OnUnauthorizedRequest() {
	r, err := http.NewRequest(http.MethodGet, ts.fullPath("not needed"), nil)
	ts.Require().NoError(err)

	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusUnauthorized, response.StatusCode)
}

func (ts *userHandlerTestSuite) TestPostUserAddsUserToStorage() {
	ts.userStoreMock.On("Create", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*model.User")).Return(nil)

	rawUser, err := json.Marshal(postUserRequest{
		Username:    "hunter",
		DisplayName: "Hunter",
		Password:    "hunter2",
	})
	ts.Require().NoError(err)
	r, err := http.NewRequest(http.MethodPost, ts.fullPath(""), bytes.NewReader(rawUser))
	ts.Require().NoError(err)
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusCreated, response.StatusCode)
}

func (ts *userHandlerTestSuite) TestPostUserRequiresPassword() {
	rawUser, err := json.Marshal(postUserRequest{
		Username:    "hunter",
		DisplayName: "Hunter",
	})
	ts.Require().NoError(err)
	r, err := http.NewRequest(http.MethodPost, ts.fullPath(""), bytes.NewReader(rawUser))
	ts.Require().NoError(err)
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusBadRequest, response.StatusCode)
}

func (ts *userHandlerTestSuite) TestPostUser400OnAlreadyExistingUser() {
	ts.userStoreMock.On("Create", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*model.User")).Return(storage.ErrEntityAlreadyExists)

	rawUser, err := json.Marshal(postUserRequest{
		Username:    "hunter",
		DisplayName: "Hunter",
		Password:    "hunter2",
	})
	ts.Require().NoError(err)
	r, err := http.NewRequest(http.MethodPost, ts.fullPath(""), bytes.NewReader(rawUser))
	ts.Require().NoError(err)
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusBadRequest, response.StatusCode)
}

func (ts *userHandlerTestSuite) TestPostLoginReturnsAuthToken() {
	clearPassword := "hunter2"
	encodedPassword, err := auth.GenerateEncodedPassword(clearPassword)
	ts.Require().NoError(err)
	user := &model.User{
		Username:       "hunter",
		Permission:     model.PermissionLevelAdmin,
		HashedPassword: encodedPassword,
	}
	ts.userStoreMock.On("Get", mock.AnythingOfType("*context.valueCtx"), user.Username).Return(user, nil)

	rawUser, err := json.Marshal(postLoginRequest{
		Username: "hunter",
		Password: clearPassword,
	})
	ts.Require().NoError(err)
	r, err := http.NewRequest(http.MethodPost, ts.fullPath("login"), bytes.NewReader(rawUser))
	ts.Require().NoError(err)
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusOK, response.StatusCode)

	defer response.Body.Close()
	rawBody, err := io.ReadAll(response.Body)
	ts.Require().NoError(err)
	var loginResponse postLoginResponse
	ts.Require().NoError(json.Unmarshal(rawBody, &loginResponse))

	ts.Require().EqualValues(user.Username, loginResponse.Username)
	validToken, err := auth.VerifyToken([]byte(loginResponse.AuthToken))
	ts.Require().NoError(err)
	ts.Require().NotNil(validToken)
}

func (ts *userHandlerTestSuite) TestPostLogin404OnNonExistantUser() {
	username := "notneeded"
	ts.userStoreMock.On("Get", mock.AnythingOfType("*context.valueCtx"), username).Return(nil, storage.ErrEntityNotFound)

	rawUser, err := json.Marshal(postLoginRequest{
		Username: username,
	})
	ts.Require().NoError(err)
	r, err := http.NewRequest(http.MethodPost, ts.fullPath("login"), bytes.NewReader(rawUser))
	ts.Require().NoError(err)

	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusNotFound, response.StatusCode)
}

func (ts *userHandlerTestSuite) TestPostLogin400OnPasswordMismatch() {
	clearPassword := "hunter2"
	encodedPassword, err := auth.GenerateEncodedPassword(clearPassword)
	ts.Require().NoError(err)
	user := &model.User{
		Username:       "hunter",
		Permission:     model.PermissionLevelAdmin,
		HashedPassword: encodedPassword,
	}
	ts.userStoreMock.On("Get", mock.AnythingOfType("*context.valueCtx"), user.Username).Return(user, nil)

	rawUser, err := json.Marshal(postLoginRequest{
		Username: "hunter",
		Password: "bad password",
	})
	ts.Require().NoError(err)
	r, err := http.NewRequest(http.MethodPost, ts.fullPath("login"), bytes.NewReader(rawUser))
	ts.Require().NoError(err)
	response, err := ts.client.Do(r)
	ts.Require().NoError(err)
	ts.Require().EqualValues(http.StatusBadRequest, response.StatusCode)
}

func (ts *userHandlerTestSuite) fullPath(path string) string {
	return fmt.Sprintf("%s/api/v1/users/%s", ts.server.URL, path)
}

func (ts *userHandlerTestSuite) makeAuthenticatedRequest(method string, path string, user *model.User) *http.Request {
	return makeAuthenticatedRequest(ts.T(), method, ts.fullPath(path), nil, user)
}
