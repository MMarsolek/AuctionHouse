package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MMarsolek/AuctionHouse/log"
	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/MMarsolek/AuctionHouse/server/controller/auth"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingFieldsAddsFieldsToLogs(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	ctx := log.WithLogger(context.Background(), zap.New(core).Sugar())
	runMiddlewareHandlerContext(ctx, LoggingFields, func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.Context(), "some log", "hello", "world")
	})

	require.EqualValues(t, 2, logs.Len())

	for _, log := range logs.All() {
		fieldsSeen := map[string]bool{
			"route":         false,
			"routeTemplate": false,
			"method":        false,
		}
		for _, field := range log.Context {
			fieldsSeen[field.Key] = true
		}

		for name, value := range fieldsSeen {
			require.Truef(t, value, "%s was not seen when it was expected to be", name)
		}
	}
}

func TestPanicHandlerRecoversFromAnyValuePanic(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("I am panicking!")
	})
	require.Panics(t, func() {
		handler.ServeHTTP(nil, nil)
	})

	require.NotPanics(t, func() {
		runMiddlewareHandler(PanicHandler, handler)
	})
}

func TestPanicHandlerRecoversFromErrorPanic(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(errors.New("I AM ERROR"))
	})
	require.Panics(t, func() {
		handler.ServeHTTP(nil, nil)
	})

	require.NotPanics(t, func() {
		runMiddlewareHandler(PanicHandler, handler)
	})
}

func TestTrailingSlashRemovesSlashFromPath(t *testing.T) {
	runMiddlewareHandler(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			require.True(t, strings.HasSuffix(r.URL.Path, "/"), "url should have a / in the suffix")
			RemoveTrailingSlash(next).ServeHTTP(rw, r)
		})
	}, func(w http.ResponseWriter, r *http.Request) {
		require.True(t, !strings.HasSuffix(r.URL.Path, "/"), "url should not have a / in the suffix")
	})
}

func TestVerifyAuthTokenWritesUnauthorizedOnNoToken(t *testing.T) {
	w := runMiddlewareHandler(VerifyAuthToken, func(w http.ResponseWriter, r *http.Request) {
		require.FailNow(t, "we should not get here")
	})

	require.EqualValues(t, http.StatusUnauthorized, w.Code)
}

func TestVerifyAuthTokenWritesUnauthorizedOnInvalidToken(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/path", nil)
	request.Header.Add("Authorization", "Bearer sometoken")

	w := httptest.NewRecorder()
	VerifyAuthToken(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		require.Fail(t, "we should not get here")
	})).ServeHTTP(w, request)

	require.EqualValues(t, http.StatusUnauthorized, w.Code)
}

func TestVerifyAuthTokenCallsHandlerOnValidToken(t *testing.T) {
	user := &model.User{
		Username:   "hunter",
		Permission: model.PermissionLevelAdmin,
	}
	rawToken, err := auth.NewToken(user)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/path", nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", rawToken))

	expectedStatus := http.StatusAccepted
	w := httptest.NewRecorder()
	VerifyAuthToken(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		require.EqualValues(t, user.Username, auth.ExtractUsername(r.Context()))
		require.EqualValues(t, user.Permission, auth.ExtractPermission(r.Context()))
		rw.WriteHeader(expectedStatus)
	})).ServeHTTP(w, request)

	require.EqualValues(t, expectedStatus, w.Code)
}

func TestVerifyPermissionsWritesForbiddenOnInvalidPermissions(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/path", nil)
	request = request.WithContext(auth.WithPermission(request.Context(), model.PermissionLevelAdmin))

	expectedStatus := http.StatusForbidden
	w := httptest.NewRecorder()
	VerifyPermissions(model.PermissionLevelBidder)(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		require.FailNow(t, "we should not get here")
	})).ServeHTTP(w, request)

	require.EqualValues(t, expectedStatus, w.Code)
}

func TestVerifyPermissionsProceedsToHandlerOnValidPermissions(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/path", nil)
	request = request.WithContext(auth.WithPermission(request.Context(), model.PermissionLevelAdmin))

	expectedStatus := http.StatusAccepted
	w := httptest.NewRecorder()
	VerifyPermissions(model.PermissionLevelBidder, model.PermissionLevelAdmin)(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(expectedStatus)
	})).ServeHTTP(w, request)

	require.EqualValues(t, expectedStatus, w.Code)
}

func runMiddlewareHandler(mid mux.MiddlewareFunc, f func(w http.ResponseWriter, r *http.Request)) *httptest.ResponseRecorder {
	return runMiddlewareHandlerContext(context.Background(), mid, f)
}

func runMiddlewareHandlerContext(ctx context.Context, mid mux.MiddlewareFunc, f func(w http.ResponseWriter, r *http.Request)) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/path/", nil)
	request = request.WithContext(ctx)

	w := httptest.NewRecorder()
	mid.Middleware(http.HandlerFunc(f)).ServeHTTP(w, request)
	return w
}
