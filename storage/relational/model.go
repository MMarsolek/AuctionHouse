package relational

import (
	"context"
	"strings"
	"time"

	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type baseDBModel struct {
	ID        uint64    `bun:",pk"`
	CreatedAt time.Time `bun:",nullzero,notnull"`
	UpdatedAt time.Time `bun:",nullzero,notnull"`
}

func (model *baseDBModel) updateTime() {
	model.UpdatedAt = time.Now().UTC()
}

func (model *baseDBModel) updateCreateTime() {
	model.CreatedAt = time.Now().UTC()
}

var _ bun.BeforeAppendModelHook = (*baseDBModel)(nil)

func (*baseDBModel) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	type timeUpdater interface {
		updateTime()
		updateCreateTime()
	}

	updater, ok := query.GetModel().Value().(timeUpdater)
	if !ok {
		return nil
	}

	updater.updateTime()
	if _, ok = query.(*bun.InsertQuery); ok {
		updater.updateCreateTime()
	}

	return nil
}

// User represents the model.User as it exists in storage.
type User struct {
	baseDBModel
	Username       string                `bun:",notnull,unique"`
	DisplayName    string                `bun:",notnull"`
	HashedPassword string                `bun:",notnull"`
	Permission     model.PermissionLevel `bun:",notnull"`
}

var _ bun.AfterCreateTableHook = (*User)(nil)

func (u *User) AfterCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	return createIndex(ctx, query, (*User)(nil), "username_idx", "username")
}

// ToModel transforms the User into a model.User.
func (u *User) ToModel() *model.User {
	return &model.User{
		Username:       u.Username,
		DisplayName:    u.DisplayName,
		HashedPassword: u.HashedPassword,
		Permission:     u.Permission,
	}
}

// UserToDBModel transforms the model.User into a User.
func UserToDBModel(user *model.User) *User {
	return &User{
		Username:       user.Username,
		DisplayName:    user.DisplayName,
		HashedPassword: user.HashedPassword,
		Permission:     user.Permission,
	}
}

// AuctionItem represents the model.AuctionItem as it exists in storage.
type AuctionItem struct {
	baseDBModel
	NameID      string `bun:"name_id,notnull,unique"`
	DisplayName string `bun:",notnull"`
	ImageRef    string `bun:",notnull"`
	Description string `bun:",notnull"`
}

var _ bun.AfterCreateTableHook = (*AuctionItem)(nil)

func (ai *AuctionItem) AfterCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	return createIndex(ctx, query, (*AuctionItem)(nil), "name_id_idx", "name_id")
}

// ToModel transforms the AuctionItem into a model.AuctionItem.
func (ai *AuctionItem) ToModel() *model.AuctionItem {
	return &model.AuctionItem{
		Name:        ai.DisplayName,
		ImageRef:    ai.ImageRef,
		Description: ai.Description,
	}
}

// AuctionItemToDBModel transforms the model.AuctionItem into an AuctionItem.
func AuctionItemToDBModel(auctionItem *model.AuctionItem) *AuctionItem {
	return &AuctionItem{
		NameID:      getAuctionItemNameID(auctionItem.Name),
		DisplayName: auctionItem.Name,
		ImageRef:    auctionItem.ImageRef,
		Description: auctionItem.Description,
	}
}

func getAuctionItemNameID(name string) string {
	return strings.ToLower(name)
}

// AuctionBid represents the model.AuctionBid as it exists in storage.
type AuctionBid struct {
	baseDBModel
	BidAmount int          `bun:",notnull"`
	Bidder    *User        `bun:"rel:has-one,join:bidder_id=id"`
	Item      *AuctionItem `bun:"rel:has-one,join:item_id=id"`

	BidderID uint64 `bun:",notnull"`
	ItemID   uint64 `bun:",notnull"`
}

// ToModel transforms the AuctionBid into a model.AuctionBid.
func (ab *AuctionBid) ToModel() *model.AuctionBid {
	return &model.AuctionBid{
		BidAmount: ab.BidAmount,
		Bidder:    ab.Bidder.ToModel(),
		Item:      ab.Item.ToModel(),
	}
}

func createIndex(ctx context.Context, query *bun.CreateTableQuery, model interface{}, indexName string, columnName string) error {
	_, err := query.DB().
		NewCreateIndex().
		Model(model).
		Index(indexName).
		Column(columnName).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		return errors.Wrapf(err, "unable to create index '%s' on table", indexName)
	}

	return nil
}
