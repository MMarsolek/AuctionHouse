package controller

import (
	"net/http"

	"github.com/MMarsolek/AuctionHouse/log"
)

type handlerWithError interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}

type handlerWithErrorFunc func(w http.ResponseWriter, r *http.Request) error

func (f handlerWithErrorFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

func wrapHandler(handler handlerWithErrorFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			log.Error(r.Context(), "error running handler", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
