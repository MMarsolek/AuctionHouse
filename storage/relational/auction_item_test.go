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
)

type auctionItemClientTestSuite struct {
	suite.Suite

	ctx    context.Context
	db     bun.IDB
	client *auctionItemClient
}

func (ts *auctionItemClientTestSuite) SetupSuite() {
	rawDB, err := sql.Open("sqlite", "file::memory:?_pragma=cache%3Dshared&_pragma=foreign_keys%3Dtrue")
	ts.Require().NoError(err)

	ts.ctx = context.Background()
	ts.db = bun.NewDB(rawDB, sqlitedialect.New())
	ts.client = &auctionItemClient{baseClient{ts.db}}

	ts.Require().NoError(CreateSchema(ts.ctx, ts.db))
}

func (ts *auctionItemClientTestSuite) SetupTest() {
	_, err := ts.db.NewTruncateTable().Model(&AuctionItem{}).Exec(ts.ctx)
	ts.Require().NoError(err)
}

func TestAuctionItemClient(t *testing.T) {
	suite.Run(t, new(auctionItemClientTestSuite))
}

func (ts *auctionItemClientTestSuite) TestCreateDoesNotErrorOnValidInput() {
	item := model.AuctionItem{
		Name:        "foo",
		ImageRef:    "https://notreal.com",
		Description: "my description",
	}
	ts.Require().NoError(ts.client.Create(ts.ctx, &item))
}

func (ts *auctionItemClientTestSuite) TestCreateAllowsMultipleItems() {
	firstItem := model.AuctionItem{
		Name:        "foo",
		ImageRef:    "whatever",
		Description: "some text here",
	}

	secondItem := model.AuctionItem{
		Name:        "bar",
		ImageRef:    "whatever also",
		Description: "some more text here",
	}

	ts.Require().NoError(ts.client.Create(ts.ctx, &firstItem))
	ts.Require().NoError(ts.client.Create(ts.ctx, &secondItem))
}

func (ts *auctionItemClientTestSuite) TestCreateFailsOnDuplicateItems() {
	firstItem := model.AuctionItem{
		Name:        "foo",
		ImageRef:    "whatever",
		Description: "some text here",
	}

	secondItem := model.AuctionItem{
		Name:        "foo",
		ImageRef:    "whatever also",
		Description: "some more text here",
	}

	ts.Require().NoError(ts.client.Create(ts.ctx, &firstItem))
	err := ts.client.Create(ts.ctx, &secondItem)
	ts.Require().Error(err)
	ts.Require().ErrorIs(err, storage.ErrEntityAlreadyExists)
}

func (ts *auctionItemClientTestSuite) TestGetRetrievesItem() {
	item := model.AuctionItem{
		Name:        "foo",
		ImageRef:    "https://notreal.com",
		Description: "my description",
	}
	ts.Require().NoError(ts.client.Create(ts.ctx, &item))

	retrievedItem, err := ts.client.Get(ts.ctx, item.Name)
	ts.Require().NoError(err)
	ts.Require().EqualValues(item.Description, retrievedItem.Description)
}

func (ts *auctionItemClientTestSuite) TestGetReturnsErrorWhenItemDoesNotExist() {
	_, err := ts.client.Get(ts.ctx, "does not exist")
	ts.Require().ErrorIs(err, storage.ErrEntityNotFound)
}

func (ts *auctionItemClientTestSuite) TestGetAllRetrievesAllItems() {
	items := []*model.AuctionItem{
		{
			Name:        "foo",
			ImageRef:    "https://notreal2.com",
			Description: "my description",
		},
		{
			Name:        "bar",
			ImageRef:    "https://notreal1.com",
			Description: "my other description",
		},
	}

	for _, item := range items {
		ts.Require().NoError(ts.client.Create(ts.ctx, item))
	}

	allItems, err := ts.client.GetAll(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Len(allItems, len(items))

	for i, item := range items {
		ts.Require().EqualValues(item.Name, allItems[i].Name)
		ts.Require().EqualValues(item.ImageRef, allItems[i].ImageRef)
		ts.Require().EqualValues(item.Description, allItems[i].Description)
	}
}

func (ts *auctionItemClientTestSuite) TestGetAllReturnsEmptyArrayWhenNoItems() {
	allItems, err := ts.client.GetAll(ts.ctx)
	ts.Require().NoError(err)
	ts.Require().Empty(allItems)
}

func (ts *auctionItemClientTestSuite) TestDeleteRemovesItem() {
	item := model.AuctionItem{
		Name:        "foo",
		ImageRef:    "https://notreal.com",
		Description: "my description",
	}
	ts.Require().NoError(ts.client.Create(ts.ctx, &item))
	ts.Require().NoError(ts.client.Delete(ts.ctx, item.Name))

	_, err := ts.client.Get(ts.ctx, item.Name)
	ts.Require().ErrorIs(err, storage.ErrEntityNotFound)
}

func (ts *auctionItemClientTestSuite) TestDeleteErrorsWhenItemDoesNotExist() {
	ts.Require().ErrorIs(ts.client.Delete(ts.ctx, "does not exist"), storage.ErrEntityNotFound)
}

func (ts *auctionItemClientTestSuite) TestUpdateModifiesNonZeroFields() {
	item := model.AuctionItem{
		Name:        "foo",
		ImageRef:    "https://notreal.com",
		Description: "my description",
	}
	newDescription := "updated description"
	ts.Require().NoError(ts.client.Create(ts.ctx, &item))

	ts.Require().NoError(ts.client.Update(ts.ctx, &model.AuctionItem{
		Name:        item.Name,
		Description: newDescription,
	}))

	retrievedItem, err := ts.client.Get(ts.ctx, item.Name)
	ts.Require().NoError(err)
	ts.Require().EqualValues(newDescription, retrievedItem.Description)
	ts.Require().EqualValues(item.ImageRef, retrievedItem.ImageRef)
}
