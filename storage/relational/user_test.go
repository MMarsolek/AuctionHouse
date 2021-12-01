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

type userClientTestSuite struct {
	suite.Suite

	ctx    context.Context
	db     bun.IDB
	client *userClient
}

func (ts *userClientTestSuite) SetupSuite() {
	rawDB, err := sql.Open("sqlite", "file::memory:?_pragma=cache%3Dshared&_pragma=foreign_keys%3Dtrue")
	ts.Require().NoError(err)

	ts.ctx = context.Background()
	ts.db = bun.NewDB(rawDB, sqlitedialect.New())
	ts.client = &userClient{baseClient{ts.db}}

	ts.Require().NoError(CreateSchema(ts.ctx, ts.db))
}

func (ts *userClientTestSuite) SetupTest() {
	_, err := ts.db.NewTruncateTable().Model(&User{}).Exec(ts.ctx)
	ts.Require().NoError(err)
}

func TestUserClient(t *testing.T) {
	suite.Run(t, new(userClientTestSuite))
}

func (ts *userClientTestSuite) TestCreateDoesNotErrorOnValidInput() {
	user := model.User{
		Username:       "foo",
		DisplayName:    "bar",
		HashedPassword: "1234",
		Permission:     model.PermissionLevelAdmin,
	}
	ts.Require().NoError(ts.client.Create(ts.ctx, &user))
}

func (ts *userClientTestSuite) TestCreateAllowsMultipleUsers() {
	firstUser := model.User{
		Username:       "foo",
		DisplayName:    "Foo",
		HashedPassword: "1234",
		Permission:     model.PermissionLevelAdmin,
	}

	secondUser := model.User{
		Username:       "bar",
		DisplayName:    "Bar",
		HashedPassword: "4321",
		Permission:     model.PermissionLevelAdmin,
	}

	ts.Require().NoError(ts.client.Create(ts.ctx, &firstUser))
	ts.Require().NoError(ts.client.Create(ts.ctx, &secondUser))
}

func (ts *userClientTestSuite) TestCreateFailsOnDuplicateUsers() {
	firstUser := model.User{
		Username:       "foo",
		DisplayName:    "Foo",
		HashedPassword: "1234",
		Permission:     model.PermissionLevelAdmin,
	}

	secondUser := model.User{
		Username:       "foo",
		DisplayName:    "Bar",
		HashedPassword: "4321",
		Permission:     model.PermissionLevelAdmin,
	}

	ts.Require().NoError(ts.client.Create(ts.ctx, &firstUser))
	err := ts.client.Create(ts.ctx, &secondUser)
	ts.Require().Error(err)
	ts.Require().ErrorIs(err, storage.ErrEntityAlreadyExists)
}

func (ts *userClientTestSuite) TestGetRetrievesUser() {
	user := model.User{
		Username:       "foo",
		DisplayName:    "bar",
		HashedPassword: "1234",
		Permission:     model.PermissionLevelAdmin,
	}
	ts.Require().NoError(ts.client.Create(ts.ctx, &user))

	retrievedUser, err := ts.client.Get(ts.ctx, user.Username)
	ts.Require().NoError(err)
	ts.Require().EqualValues(user.DisplayName, retrievedUser.DisplayName)
}

func (ts *userClientTestSuite) TestGetReturnsErrorWhenUserDoesNotExist() {
	_, err := ts.client.Get(ts.ctx, "does not exist")
	ts.Require().ErrorIs(err, storage.ErrEntityNotFound)
}

func (ts *userClientTestSuite) TestDeleteRemovesUser() {
	user := model.User{
		Username:       "foo",
		DisplayName:    "bar",
		HashedPassword: "1234",
		Permission:     model.PermissionLevelAdmin,
	}
	ts.Require().NoError(ts.client.Create(ts.ctx, &user))
	ts.Require().NoError(ts.client.Delete(ts.ctx, user.Username))

	_, err := ts.client.Get(ts.ctx, user.Username)
	ts.Require().ErrorIs(err, storage.ErrEntityNotFound)
}

func (ts *userClientTestSuite) TestDeleteErrorsWhenUserDoesNotExist() {
	ts.Require().ErrorIs(ts.client.Delete(ts.ctx, "does not exist"), storage.ErrEntityNotFound)
}

func (ts *userClientTestSuite) TestUpdateModifiesNonZeroFields() {
	user := model.User{
		Username:       "foo",
		DisplayName:    "bar",
		HashedPassword: "1234",
		Permission:     model.PermissionLevelAdmin,
	}
	newDisplayName := "updated display name"
	ts.Require().NoError(ts.client.Create(ts.ctx, &user))

	ts.Require().NoError(ts.client.Update(ts.ctx, &model.User{
		Username:    user.Username,
		DisplayName: newDisplayName,
	}))

	retrievedUser, err := ts.client.Get(ts.ctx, user.Username)
	ts.Require().NoError(err)
	ts.Require().EqualValues(newDisplayName, retrievedUser.DisplayName)
	ts.Require().EqualValues(user.HashedPassword, retrievedUser.HashedPassword)
}
