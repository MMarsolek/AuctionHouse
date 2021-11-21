package relational

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/MMarsolek/AuctionHouse/storage"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type testModel struct {
	baseDBModel
	OtherPrimaryKey string `bun:",unique"`
	Value           int
}

type baseClientTestSuite struct {
	suite.Suite

	ctx    context.Context
	db     bun.IDB
	client *baseClient
}

func (ts *baseClientTestSuite) SetupSuite() {
	rawDB, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	ts.Require().NoError(err)

	ts.ctx = context.Background()
	ts.db = bun.NewDB(rawDB, sqlitedialect.New())
	ts.client = &baseClient{ts.db}

	_, err = ts.db.NewCreateTable().Model(&testModel{}).Exec(ts.ctx)
	ts.Require().NoError(err)
}

func (ts *baseClientTestSuite) SetupTest() {
	_, err := ts.db.NewTruncateTable().Model(&testModel{}).Exec(ts.ctx)
	ts.Require().NoError(err)
}

func TestBaseClient(t *testing.T) {
	suite.Run(t, new(baseClientTestSuite))
}

func (ts *baseClientTestSuite) TestCreateDoesNotErrorOnNewEntry() {
	model := testModel{
		OtherPrimaryKey: "first",
		Value:           10,
	}
	ts.Require().NoError(ts.client.create(ts.ctx, &model))
}

func (ts *baseClientTestSuite) TestCreateErrorsWhenNonPointerPassed() {
	ts.Require().Error(ts.client.create(ts.ctx, testModel{}))
}

func (ts *baseClientTestSuite) TestGetReturnsExistingItem() {
	model := testModel{
		OtherPrimaryKey: "first",
		Value:           5,
	}

	ts.Require().NoError(ts.client.create(ts.ctx, &model))
	var retrievedModel testModel
	ts.Require().NoError(ts.client.get(ts.ctx, &retrievedModel, "other_primary_key", model.OtherPrimaryKey))

	ts.Require().EqualValues(model.ID, retrievedModel.ID)
	ts.Require().EqualValues(model.OtherPrimaryKey, retrievedModel.OtherPrimaryKey)
	ts.Require().EqualValues(model.Value, retrievedModel.Value)
	ts.Require().WithinDuration(model.CreatedAt, retrievedModel.CreatedAt, time.Second)
	ts.Require().WithinDuration(model.UpdatedAt, retrievedModel.UpdatedAt, time.Second)
}

func (ts *baseClientTestSuite) TestGetReturnsErrEntityNotFoundWhenNotFound() {
	var retrievedModel testModel
	err := ts.client.get(ts.ctx, &retrievedModel, "other_primary_key", "not exists")
	ts.Require().Error(err)
	ts.Require().ErrorIs(err, storage.ErrEntityNotFound)
}

func (ts *baseClientTestSuite) TestGetErrorsWhenNonPointerPassed() {
	ts.Require().Error(ts.client.get(ts.ctx, testModel{}, "", ""))
}

func (ts *baseClientTestSuite) TestGetAllReturnsAllItems() {
	models := []*testModel{
		{
			OtherPrimaryKey: "first",
			Value:           10,
		},
		{
			OtherPrimaryKey: "second",
			Value:           -5,
		},
	}

	for _, model := range models {
		ts.Require().NoError(ts.client.create(ts.ctx, model))
	}

	var retrievedModels []*testModel
	ts.Require().NoError(ts.client.getAll(ts.ctx, &retrievedModels))
	ts.Require().Len(retrievedModels, len(models))
	ts.Require().EqualValues(models[0].OtherPrimaryKey, retrievedModels[0].OtherPrimaryKey)
	ts.Require().EqualValues(models[0].Value, retrievedModels[0].Value)
	ts.Require().EqualValues(models[1].OtherPrimaryKey, retrievedModels[1].OtherPrimaryKey)
	ts.Require().EqualValues(models[1].Value, retrievedModels[1].Value)
}

func (ts *baseClientTestSuite) TestGetAllDoesNotErrorOnEmpty() {
	var retrievedModels []*testModel
	ts.Require().NoError(ts.client.getAll(ts.ctx, &retrievedModels))
	ts.Require().Empty(retrievedModels)
}

func (ts *baseClientTestSuite) TestGetAllErrorsWhenNonPointerPassed() {
	ts.Require().Error(ts.client.getAll(ts.ctx, nil))
}

func (ts *baseClientTestSuite) TestDeleteReturnsErrEntityNotFoundWhenNotFound() {
	var retrievedModel testModel
	err := ts.client.delete(ts.ctx, &retrievedModel, "other_primary_key", "not exists")
	ts.Require().Error(err)
	ts.Require().ErrorIs(err, storage.ErrEntityNotFound)
}

func (ts *baseClientTestSuite) TestDeleteRemovesEntityFromTable() {
	model := testModel{
		OtherPrimaryKey: "second",
		Value:           -1,
	}
	ts.Require().NoError(ts.client.create(ts.ctx, &model))
	ts.Require().NoError(ts.client.delete(ts.ctx, &model, "id", model.ID))
	var retrievedModel testModel
	err := ts.client.get(ts.ctx, &retrievedModel, "id", model.ID)
	ts.Require().Error(err)
	ts.Require().ErrorIs(err, storage.ErrEntityNotFound)
}

func (ts *baseClientTestSuite) TestDeleteErrorsWhenNonPointerPassed() {
	ts.Require().Error(ts.client.delete(ts.ctx, testModel{}, "", ""))
}

func (ts *baseClientTestSuite) TestUpdateModifiesSelectedColumns() {
	model := testModel{
		OtherPrimaryKey: "unchanged",
		Value:           42,
	}
	ts.Require().NoError(ts.client.create(ts.ctx, &model))

	updatedModel := testModel{
		baseDBModel: baseDBModel{
			ID: model.ID,
		},
		OtherPrimaryKey: "changed",
		Value:           10,
	}

	ts.Require().NoError(ts.client.update(ts.ctx, &updatedModel, "id", model.ID, "value"))
	var retrievedModel testModel
	ts.Require().NoError(ts.client.get(ts.ctx, &retrievedModel, "id", model.ID))

	ts.Require().EqualValues(model.ID, retrievedModel.ID)
	ts.Require().EqualValues(updatedModel.OtherPrimaryKey, retrievedModel.OtherPrimaryKey)
	ts.Require().EqualValues(updatedModel.Value, retrievedModel.Value)
	ts.Require().NotEqualValues(model.UpdatedAt, retrievedModel.UpdatedAt)
}

func (ts *baseClientTestSuite) TestUpdateErrorsWhenNonPointerPassed() {
	ts.Require().Error(ts.client.update(ts.ctx, testModel{}, "", ""))
}
