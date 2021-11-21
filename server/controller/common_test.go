package controller

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/MMarsolek/AuctionHouse/server/controller/auth"
	"github.com/MMarsolek/AuctionHouse/server/controller/middleware"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func newTestServer() (*httptest.Server, *mux.Router) {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	return httptest.NewServer(middleware.RemoveTrailingSlash(router)), router
}

func makeAuthenticatedRequest(t *testing.T, method string, fullPath string, body io.Reader, user *model.User) *http.Request {
	r, err := http.NewRequest(method, fullPath, body)
	require.NoError(t, err)

	token, err := auth.NewToken(user)
	require.NoError(t, err)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return r
}
