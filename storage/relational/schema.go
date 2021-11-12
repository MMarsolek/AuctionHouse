package relational

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func CreateSchema(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().
		Model((*User)(nil)).
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "unable to create users table")
	}

	_, err = db.NewCreateTable().
		Model((*AuctionItem)(nil)).
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "unable to create auction_items table")
	}

	_, err = db.NewCreateTable().
		Model((*AuctionBid)(nil)).
		ForeignKey(`("bidder_id") REFERENCES "users" ("id")`).
		ForeignKey(`("item_id") REFERENCES "auction_items" ("id")`).
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "unable to create auction_bids table")
	}

	return nil
}
