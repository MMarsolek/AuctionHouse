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

type UserClient interface {
	Get(ctx context.Context, username string) (*model.User, error)
	Delete(ctx context.Context, username string) error
	Update(ctx context.Context, user *model.User) error
	Create(ctx context.Context, user *model.User) error
}

type AuctionItemClient interface {
	Get(ctx context.Context, name string) (*model.AuctionItem, error)
	Delete(ctx context.Context, name string) error
	Update(ctx context.Context, item *model.AuctionItem) error
	Create(ctx context.Context, item *model.AuctionItem) error
}

type AuctionBidClient interface {
	GetHighestBid(ctx context.Context, item *model.AuctionItem) (*model.AuctionBid, error)
	GetAllHighestBids(ctx context.Context) ([]*model.AuctionBid, error)
	PlaceBid(ctx context.Context, user *model.User, item *model.AuctionItem, amount int) (*model.AuctionBid, error)
}
