package storage

import (
	"context"
	"errors"

	"github.com/MMarsolek/AuctionHouse/model"
)

var (
	ErrEntityNotFound      = errors.New("entity not found")
	ErrEntityAlreadyExists = errors.New("entity already exists")
)

type UserClient interface {
	Get(ctx context.Context, username string) (*model.User, error)
	Delete(ctx context.Context, username string) error
	Update(ctx context.Context, user *model.User) error
	Create(ctx context.Context, user *model.User) error
}
