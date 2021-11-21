package relational

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func CreateSchema(ctx context.Context, db bun.IDB) error {
	_, err := db.NewCreateTable().
		Model((*User)(nil)).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "unable to create users table")
	}

	_, err = db.NewCreateTable().
		Model((*AuctionItem)(nil)).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "unable to create auction_items table")
	}

	_, err = db.NewCreateTable().
		Model((*AuctionBid)(nil)).
		IfNotExists().
		ForeignKey(`("bidder_id") REFERENCES "users" ("id")`).
		ForeignKey(`("item_id") REFERENCES "auction_items" ("id")`).
		Exec(ctx)

	if err != nil {
		return errors.Wrap(err, "unable to create auction_bids table")
	}

	_, err = db.ExecContext(ctx, `
	CREATE TRIGGER IF NOT EXISTS highest_value_check
	BEFORE INSERT ON auction_bids
	BEGIN
		SELECT RAISE(FAIL, "cannot bid lower")
		FROM auction_bids
		WHERE item_id = NEW.item_id
		GROUP BY bidder_id, item_id
		HAVING MAX(bid_amount) >= NEW.bid_amount;
	END;`)

	if err != nil {
		return errors.Wrap(err, "unable to create trigger on auction_bids table")
	}

	return nil
}
