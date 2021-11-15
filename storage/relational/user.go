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

func NewUserClient(db bun.IDB) storage.UserClient {
	return &userClient{
		baseClient: baseClient{
			db: db,
		},
	}
}

func (uc *userClient) Get(ctx context.Context, username string) (*model.User, error) {
	var user User
	err := uc.baseClient.Get(ctx, &user, "username", username)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get user with username '%s'", username)
	}
	return user.ToModel(), nil
}

func (uc *userClient) Delete(ctx context.Context, username string) error {
	err := uc.baseClient.Delete(ctx, (*User)(nil), "username", username)
	if err != nil {
		return errors.Wrapf(err, "unable to delete user with username '%s'", username)
	}
	return nil
}

func (uc *userClient) Update(ctx context.Context, user *model.User) error {
	dbModel := UserToDBModel(user)
	err := uc.baseClient.Update(ctx, dbModel, "username", user.Username)
	if err != nil {
		return errors.Wrapf(err, "unable to update user %s", user.Username)
	}
	return nil
}

func (uc *userClient) Create(ctx context.Context, user *model.User) error {
	dbModel := UserToDBModel(user)
	err := uc.baseClient.Create(ctx, dbModel)
	if err != nil {
		return errors.Wrapf(err, "unable to create user %s", user.Username)
	}

	*user = *dbModel.ToModel()
	return nil
}
