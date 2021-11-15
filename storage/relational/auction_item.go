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

func NewAuctionItemClient(db bun.IDB) storage.AuctionItemClient {
	return &auctionItemClient{
		baseClient: baseClient{
			db: db,
		},
	}
}

func (ac *auctionItemClient) Get(ctx context.Context, name string) (*model.AuctionItem, error) {
	var item AuctionItem
	nameID := getAuctionItemNameID(name)
	err := ac.baseClient.Get(ctx, &item, "name_id", nameID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get auction item with name '%s'", nameID)
	}
	return item.ToModel(), nil
}

func (ac *auctionItemClient) Delete(ctx context.Context, name string) error {
	nameID := getAuctionItemNameID(name)
	err := ac.baseClient.Delete(ctx, (*AuctionItem)(nil), "name_id", nameID)
	if err != nil {
		return errors.Wrapf(err, "unable to delete auction item with name '%s'", nameID)
	}
	return nil
}

func (ac *auctionItemClient) Update(ctx context.Context, item *model.AuctionItem) error {
	dbModel := AuctionItemToDBModel(item)
	nameID := getAuctionItemNameID(item.Name)
	err := ac.baseClient.Update(ctx, dbModel, "name_id", nameID)
	if err != nil {
		return errors.Wrapf(err, "unable to update auction item %s", nameID)
	}
	return nil
}

func (ac *auctionItemClient) Create(ctx context.Context, item *model.AuctionItem) error {
	dbModel := AuctionItemToDBModel(item)
	nameID := getAuctionItemNameID(item.Name)
	err := ac.baseClient.Create(ctx, dbModel)
	if err != nil {
		return errors.Wrapf(err, "unable to create auction item %s", nameID)
	}

	*item = *dbModel.ToModel()
	return nil
}
