// Package classification AuctionHouse API.
//
// The purpose of this application is to provide the client with functionality to define users, items for auctions,
// and ways to place bids on those items.
//
//  Terms Of Service:
//    There are no terms of service at the moment. Use at your own risk, we take no responsibility.
//
//  Schemes: http
//  Host: localhost
//  BasePath: /api/v1
//  Version: 0.1.0
//  License: MIT http://opensource.org/licenses/MIT
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
//
//  SecurityDefinitions:
//  api_key:
//    type: apiKey
//    name: KEY
//    in: header
//
// swagger:meta
package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/MMarsolek/AuctionHouse/server/controller"
	"github.com/MMarsolek/AuctionHouse/server/controller/middleware"
	"github.com/MMarsolek/AuctionHouse/server/controller/ws"
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
	router.Handle("/", fs)

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

	websocketHandler := ws.NewHandler(userClient, itemClient, bidClient)
	websocketHandler.RegisterRoutes(rootRouter)
}
