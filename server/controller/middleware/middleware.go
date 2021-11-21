package middleware

import (
	"net/http"
	"strings"

	"github.com/MMarsolek/AuctionHouse/log"
	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/MMarsolek/AuctionHouse/server/controller/auth"
	"github.com/gorilla/mux"
)

// LoggingFields adds handler specific fields to the logging object.
func LoggingFields(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			routeTemplate = "unknown"
			err           error
		)
		if currentRoute := mux.CurrentRoute(r); currentRoute != nil {
			routeTemplate, err = mux.CurrentRoute(r).GetPathTemplate()
			if err != nil {
				log.Error(r.Context(), "could not get route template")
			}
		}

		r = r.WithContext(log.WithFields(r.Context(),
			"route", r.URL.Path,
			"routeTemplate", routeTemplate,
			"method", r.Method,
		))

		log.Info(r.Context(), "web request")
		next.ServeHTTP(w, r)
	})
}

// PanicHandler recovers any unhandled panics from the handler.
func PanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if raw := recover(); raw != nil {
				if err, ok := raw.(error); ok {
					log.Error(r.Context(), "unhandled panic error", "err", err)
				} else {
					log.Error(r.Context(), "unhandled panic with non-error", "err", raw)
				}
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// RemoveTrailingSlash removes the trailing slash from the request path.
func RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		r.URL.RawPath = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(w, r)
	})
}

// VerifyAuthToken prevents moving to the next handler if a token is not supplied or it's invalid.
func VerifyAuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authValue := r.Header.Get(http.CanonicalHeaderKey("Authorization"))
		parts := strings.Split(authValue, " ")
		if len(parts) != 2 || !strings.EqualFold("Bearer", parts[0]) {
			log.Info(r.Context(), "malformed request authorization", "authorization", authValue)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := parts[1]
		payload, err := auth.VerifyToken([]byte(token))
		if err != nil {
			log.Info(r.Context(), "invalid token", "token", token)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = r.WithContext(log.WithFields(r.Context(), "username", payload.Username, "permission", payload.Permission))
		r = r.WithContext(auth.WithUsername(r.Context(), payload.Username))
		r = r.WithContext(auth.WithPermission(r.Context(), payload.Permission))
		next.ServeHTTP(w, r)
	})
}

// VerifyPermissions prevents moving to the next handler if the context does not contain the provided permissions.
func VerifyPermissions(permissions ...model.PermissionLevel) mux.MiddlewareFunc {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userPermission := auth.ExtractPermission(r.Context())
			found := false
			for _, permission := range permissions {
				if permission == userPermission {
					found = true
					break
				}
			}

			if !found {
				log.Info(r.Context(), "insufficent permissions", "requiredPermissions", permissions, "suppliedPermission", userPermission)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	})
}
