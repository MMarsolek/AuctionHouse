package relational

import (
	"context"
	"database/sql"

	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/MMarsolek/AuctionHouse/storage"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"modernc.org/sqlite"
)

type auctionBidClient struct {
	baseClient
}

// NewAuctionBidClient returns an object that can perform various operations on model.AuctionBids.
func NewAuctionBidClient(db bun.IDB) storage.AuctionBidClient {
	return &auctionBidClient{
		baseClient{
			db: db,
		},
	}
}

// GetHighestBid gets the highest bid for the specified item. This will return storage.ErrEntityNotFound if the item
// does not have a bid.
func (bc *auctionBidClient) GetHighestBid(ctx context.Context, item *model.AuctionItem) (*model.AuctionBid, error) {
	var bid AuctionBid
	err := bc.db.NewSelect().
		Model(&bid).
		Relation("Bidder").
		Relation("Item").
		Where("item_id = (?)", bc.getRelatedItemQuery(item)).
		Group("item_id").
		Having("MAX(bid_amount)").
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(storage.ErrEntityNotFound, "item does not have a bid")
		}
		return nil, errors.Wrap(err, "retrieving highest bid")
	}
	return bid.ToModel(), nil
}

// GetAllHighestBids gets the highest bid for all items in storage.
func (bc *auctionBidClient) GetAllHighestBids(ctx context.Context) ([]*model.AuctionBid, error) {
	var bids []*AuctionBid
	err := bc.db.NewSelect().
		Model(&bids).
		Relation("Bidder").
		Relation("Item").
		Group("item_id").
		Having("MAX(bid_amount)").
		Scan(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "retrieving all highest bids")
	}

	result := make([]*model.AuctionBid, len(bids))
	for i, bid := range bids {
		result[i] = bid.ToModel()
	}
	return result, nil
}

// PlaceBid creates a new bid by the specified user for the specified item. This will return storage.ErrBidTooLow if the
// specified amount is lower than another bid in storage.
func (bc *auctionBidClient) PlaceBid(ctx context.Context, user *model.User, item *model.AuctionItem, amount int) (*model.AuctionBid, error) {
	bid := &AuctionBid{
		BidAmount: amount,
		Bidder:    UserToDBModel(user),
		Item:      AuctionItemToDBModel(item),
	}

	_, err := bc.db.NewInsert().
		Model(bid).
		Value("bidder_id", "(?)", bc.getRelatedUserQuery(user)).
		Value("item_id", "(?)", bc.getRelatedItemQuery(item)).
		Exec(ctx)

	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == 1811 {
			return nil, errors.Wrap(storage.ErrBidTooLow, "unable to insert new auction bid")
		}
		return nil, errors.Wrap(err, "inserting new auction bid")
	}
	return nil, nil
}

func (bc *auctionBidClient) getRelatedItemQuery(item *model.AuctionItem) *bun.SelectQuery {
	return bc.db.NewSelect().
		Model(item).
		Column("id").
		Where("name_id = ?", getAuctionItemNameID(item.Name))
}

func (bc *auctionBidClient) getRelatedUserQuery(user *model.User) *bun.SelectQuery {
	return bc.db.NewSelect().
		Model(user).
		Column("id").
		Where("username = ?", user.Username)
}
