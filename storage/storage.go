package storage

import (
	"context"
	"errors"

	"github.com/MMarsolek/AuctionHouse/model"
)

var (
	ErrEntityNotFound      = errors.New("entity not found")
	ErrEntityAlreadyExists = errors.New("entity already exists")
	ErrBidTooLow           = errors.New("bid is lower than current bid")
)

// UserClient defines how to store model.User objects.
//go:generate mockery --name UserClient
type UserClient interface {

	// Get retrieves the model from storage.
	Get(ctx context.Context, username string) (*model.User, error)

	// Delete does a hard remove from storage.
	Delete(ctx context.Context, username string) error

	// Update changes the non-zero fields in the supplied model.
	Update(ctx context.Context, user *model.User) error

	// Create adds a new model to storage.
	Create(ctx context.Context, user *model.User) error
}

// AuctionItemClient defines how to store model.AuctionItem objects.
//go:generate mockery --name AuctionItemClient
type AuctionItemClient interface {

	// Get retrieves the model from storage.
	Get(ctx context.Context, name string) (*model.AuctionItem, error)

	// GetAll retrieves all models from storage.
	GetAll(ctx context.Context) ([]*model.AuctionItem, error)

	// Delete does a hard remove from storage.
	Delete(ctx context.Context, name string) error

	// Update changes the non-zero fields in the supplied model.
	Update(ctx context.Context, item *model.AuctionItem) error

	// Create adds a new model to storage.
	Create(ctx context.Context, item *model.AuctionItem) error
}

// AuctionItemClient defines how to store model.AuctionBid objects.
//go:generate mockery --name AuctionBidClient
type AuctionBidClient interface {

	// GetHighestBid retrieves the highest bid for the item.
	GetHighestBid(ctx context.Context, item *model.AuctionItem) (*model.AuctionBid, error)

	// GetAllHighestBids retrieves the highest bid for every item.
	GetAllHighestBids(ctx context.Context) ([]*model.AuctionBid, error)

	// PlaceBid makes a new bid for the item by the supplied user.
	PlaceBid(ctx context.Context, user *model.User, item *model.AuctionItem, amount int) (*model.AuctionBid, error)
}
