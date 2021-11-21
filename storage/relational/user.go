package relational

import (
	"context"

	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/MMarsolek/AuctionHouse/storage"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type userClient struct {
	baseClient
}

// NewUserClient returns an object that can perform various operations on model.Users.
func NewUserClient(db bun.IDB) storage.UserClient {
	return &userClient{
		baseClient: baseClient{
			db: db,
		},
	}
}

// Get retrieves the model.User by the username. This will return storage.ErrEntityNotFound if the username is not
// found in storage.
func (uc *userClient) Get(ctx context.Context, username string) (*model.User, error) {
	var user User
	err := uc.baseClient.get(ctx, &user, "username", username)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get user with username '%s'", username)
	}
	return user.ToModel(), nil
}

// Delete removes the model.User from storage by username. This will return storage.ErrEntityNotFound if the username
// is not found in storage.
func (uc *userClient) Delete(ctx context.Context, username string) error {
	err := uc.baseClient.delete(ctx, (*User)(nil), "username", username)
	if err != nil {
		return errors.Wrapf(err, "unable to delete user with username '%s'", username)
	}
	return nil
}

// Update changes the existing user by the non-zero fields of the provided model.User object.
func (uc *userClient) Update(ctx context.Context, user *model.User) error {
	dbModel := UserToDBModel(user)
	err := uc.baseClient.update(ctx, dbModel, "username", user.Username)
	if err != nil {
		return errors.Wrapf(err, "unable to update user %s", user.Username)
	}
	return nil
}

// Create adds a new model.User to storage. This will return storage.ErrEntityAlreadyExists if the username is already
// found in storage.
func (uc *userClient) Create(ctx context.Context, user *model.User) error {
	dbModel := UserToDBModel(user)
	err := uc.baseClient.create(ctx, dbModel)
	if err != nil {
		return errors.Wrapf(err, "unable to create user %s", user.Username)
	}

	*user = *dbModel.ToModel()
	return nil
}
