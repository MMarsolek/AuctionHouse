package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/MMarsolek/AuctionHouse/server/controller"
	"github.com/MMarsolek/AuctionHouse/server/controller/middleware"
	"github.com/MMarsolek/AuctionHouse/storage"
	"github.com/gorilla/mux"
)

func NewAuctionHouseServer(
	ctx context.Context,
	userClient storage.UserClient,
	auctionItemClient storage.AuctionItemClient,
	auctionBidClient storage.AuctionBidClient,
	address string,
) *http.Server {
	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("./web"))
	router.Handle("/web", http.StripPrefix("/web", fs))

	setupControllers(router.PathPrefix("/api").Subrouter(), userClient, auctionItemClient, auctionBidClient)

	return &http.Server{
		BaseContext:  func(net.Listener) context.Context { return ctx },
		Handler:      middleware.RemoveTrailingSlash(router),
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func setupControllers(
	rootRouter *mux.Router,
	userClient storage.UserClient,
	itemClient storage.AuctionItemClient,
	bidClient storage.AuctionBidClient,
) {
	rootRouter.Use(middleware.LoggingFields)
	rootRouter.Use(middleware.PanicHandler)

	userHandler := controller.NewUserHandler(userClient)
	userHandler.RegisterRoutes(rootRouter)

	auctionHandler := controller.NewAuctionHandler(userClient, itemClient, bidClient)
	auctionHandler.RegisterRoutes(rootRouter)
}
