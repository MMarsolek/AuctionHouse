package relational

import (
	"context"
	"database/sql"
	"testing"

	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/MMarsolek/AuctionHouse/storage"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type auctionBidClientTestSuite struct {
	suite.Suite

	ctx        context.Context
	db         bun.IDB
	client     *auctionBidClient
	userClient *userClient
	itemClient *auctionItemClient
}

func (ts *auctionBidClientTestSuite) SetupSuite() {
	rawDB, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	ts.Require().NoError(err)

	ts.ctx = context.Background()
	ts.db = bun.NewDB(rawDB, sqlitedialect.New())
	ts.client = &auctionBidClient{baseClient{ts.db}}
	ts.userClient = &userClient{baseClient{ts.db}}
	ts.itemClient = &auctionItemClient{baseClient{ts.db}}

	ts.Require().NoError(CreateSchema(ts.ctx, ts.db))
}

func (ts *auctionBidClientTestSuite) SetupTest() {
	models := []interface{}{
		&AuctionBid{},
		&User{},
		&AuctionItem{},
	}
	for _, model := range models {
		_, err := ts.db.NewTruncateTable().Model(model).Exec(ts.ctx)
		ts.Require().NoError(err)
	}
}

func TestAuctionBidClient(t *testing.T) {
	suite.Run(t, new(auctionBidClientTestSuite))
}

func (ts *auctionBidClientTestSuite) TestPlaceBidAddsNewBids() {
	users, items := ts.createTestAssets()

	_, err := ts.client.PlaceBid(ts.ctx, users[0], items[0], 10)
	ts.Require().NoError(err)
}

func (ts *auctionBidClientTestSuite) TestPlaceBidDoesNotAllowSmallerBidsForItem() {
	users, items := ts.createTestAssets()

	_, err := ts.client.PlaceBid(ts.ctx, users[0], items[0], 10)
	ts.Require().NoError(err)

	_, err = ts.client.PlaceBid(ts.ctx, users[1], items[0], 20)
	ts.Require().NoError(err)

	_, err = ts.client.PlaceBid(ts.ctx, users[0], items[0], 15)
	ts.Require().Error(err)
	ts.Require().ErrorIs(err, storage.ErrBidTooLow)

	_, err = ts.client.PlaceBid(ts.ctx, users[0], items[0], 40)
	ts.Require().NoError(err)
}

func (ts *auctionBidClientTestSuite) TestGetHighestBidReturnsHighestBidForItem() {
	users, items := ts.createTestAssets()

	_, err := ts.client.PlaceBid(ts.ctx, users[0], items[0], 10)
	ts.Require().NoError(err)

	_, err = ts.client.PlaceBid(ts.ctx, users[1], items[0], 20)
	ts.Require().NoError(err)

	_, err = ts.client.PlaceBid(ts.ctx, users[0], items[1], 10)
	ts.Require().NoError(err)

	firstItemBid, err := ts.client.GetHighestBid(ts.ctx, items[0])
	ts.Require().NoError(err)
	ts.Require().EqualValues(20, firstItemBid.BidAmount)

	secondItemBid, err := ts.client.GetHighestBid(ts.ctx, items[1])
	ts.Require().NoError(err)
	ts.Require().EqualValues(10, secondItemBid.BidAmount)
}

func (ts *auctionBidClientTestSuite) TestGetHighestBidReturnsEntityNotFoundWhenNoHighestBid() {
	_, items := ts.createTestAssets()

	_, err := ts.client.GetHighestBid(ts.ctx, items[0])
	ts.Require().Error(err)
	ts.Require().ErrorIs(err, storage.ErrEntityNotFound)
}

func (ts *auctionBidClientTestSuite) TestGetAllHighestBidReturnsHighestBidForEachItem() {
	users, items := ts.createTestAssets()

	_, err := ts.client.PlaceBid(ts.ctx, users[0], items[0], 10)
	ts.Require().NoError(err)

	_, err = ts.client.PlaceBid(ts.ctx, users[1], items[0], 20)
	ts.Require().NoError(err)

	_, err = ts.client.PlaceBid(ts.ctx, users[0], items[1], 10)
	ts.Require().NoError(err)

	_, err = ts.client.PlaceBid(ts.ctx, users[2], items[1], 100)
	ts.Require().NoError(err)

	highestBids, err := ts.client.GetAllHighestBids(ts.ctx)
	ts.Require().NoError(err)

	ts.Require().Len(highestBids, 2)
	ts.Require().EqualValues(20, highestBids[0].BidAmount)
	ts.Require().EqualValues(users[1].DisplayName, highestBids[0].Bidder.DisplayName)
	ts.Require().EqualValues(100, highestBids[1].BidAmount)
	ts.Require().EqualValues(users[2].DisplayName, highestBids[1].Bidder.DisplayName)
}

func (ts *auctionBidClientTestSuite) createTestAssets() ([]*model.User, []*model.AuctionItem) {
	users := []*model.User{
		{
			Username:       "user1",
			DisplayName:    "USER 1",
			HashedPassword: "1234",
			Permission:     model.PermissionLevelBidder,
		},
		{
			Username:       "user2",
			DisplayName:    "USER 2",
			HashedPassword: "1234",
			Permission:     model.PermissionLevelBidder,
		},
		{
			Username:       "user3",
			DisplayName:    "USER 3",
			HashedPassword: "1234",
			Permission:     model.PermissionLevelBidder,
		},
	}

	for _, user := range users {
		ts.Require().NoError(ts.userClient.Create(ts.ctx, user))
	}

	items := []*model.AuctionItem{
		{
			Name:        "item1",
			ImageRef:    "some image",
			Description: "this is the first image",
		},
		{
			Name:        "item2",
			ImageRef:    "some image",
			Description: "this is the second image",
		},
	}

	for _, item := range items {
		ts.Require().NoError(ts.itemClient.Create(ts.ctx, item))
	}

	return users, items
}
