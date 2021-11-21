package relational

import (
	"context"

	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/MMarsolek/AuctionHouse/storage"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type auctionItemClient struct {
	baseClient
}

// NewAuctionItemClient returns an object that can perform various operations on model.AuctionItems.
func NewAuctionItemClient(db bun.IDB) storage.AuctionItemClient {
	return &auctionItemClient{
		baseClient: baseClient{
			db: db,
		},
	}
}

// Get retrieves the model.AuctionItem by the name. This will return storage.ErrEntityNotFound if the name is not found
// in storage.
func (ac *auctionItemClient) Get(ctx context.Context, name string) (*model.AuctionItem, error) {
	var item AuctionItem
	nameID := getAuctionItemNameID(name)
	err := ac.baseClient.get(ctx, &item, "name_id", nameID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get auction item with name '%s'", nameID)
	}
	return item.ToModel(), nil
}

// GetAll retrieves all model.AuctionItems in storage.
func (ac *auctionItemClient) GetAll(ctx context.Context) ([]*model.AuctionItem, error) {
	var dbModels []*AuctionItem
	if err := ac.baseClient.getAll(ctx, &dbModels); err != nil {
		return nil, errors.Wrapf(err, "unable to get all auction items")
	}

	result := make([]*model.AuctionItem, len(dbModels))
	for i, dbModel := range dbModels {
		result[i] = dbModel.ToModel()
	}

	return result, nil
}

// Delete removes the model.AuctionItem from storage by name. This will return storage.ErrEntityNotFound if the name is
// not found in storage.
func (ac *auctionItemClient) Delete(ctx context.Context, name string) error {
	nameID := getAuctionItemNameID(name)
	err := ac.baseClient.delete(ctx, (*AuctionItem)(nil), "name_id", nameID)
	if err != nil {
		return errors.Wrapf(err, "unable to delete auction item with name '%s'", nameID)
	}
	return nil
}

// Update changes the existing item by the non-zero fields of the provided model.AuctionItem object.
func (ac *auctionItemClient) Update(ctx context.Context, item *model.AuctionItem) error {
	dbModel := AuctionItemToDBModel(item)
	nameID := getAuctionItemNameID(item.Name)
	err := ac.baseClient.update(ctx, dbModel, "name_id", nameID)
	if err != nil {
		return errors.Wrapf(err, "unable to update auction item %s", nameID)
	}
	return nil
}

// Create adds a new model.AuctionItem to storage. This will return storage.ErrEntityAlreadyExists if the name is already
// found in storage.
func (ac *auctionItemClient) Create(ctx context.Context, item *model.AuctionItem) error {
	dbModel := AuctionItemToDBModel(item)
	nameID := getAuctionItemNameID(item.Name)
	err := ac.baseClient.create(ctx, dbModel)
	if err != nil {
		return errors.Wrapf(err, "unable to create auction item %s", nameID)
	}

	*item = *dbModel.ToModel()
	return nil
}
