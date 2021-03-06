package relational

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/MMarsolek/AuctionHouse/storage"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"modernc.org/sqlite"
)

type baseClient struct {
	db bun.IDB
}

func (bc *baseClient) get(ctx context.Context, model interface{}, primaryCol string, primaryKey interface{}) error {
	err := bc.validatePointer(model)
	if err != nil {
		return errors.Wrap(err, "unable to validate pointer")
	}
	err = bc.db.NewSelect().Model(model).Where("? = ?", bun.Ident(primaryCol), primaryKey).Scan(ctx)
	if err == sql.ErrNoRows {
		return errors.Wrapf(storage.ErrEntityNotFound, "unable to find '%s':'%s'", primaryCol, primaryKey)
	} else if err != nil {
		return errors.Wrapf(err, "unable to retrieve entity %s", primaryKey)
	}

	return nil
}

func (bc *baseClient) getAll(ctx context.Context, model interface{}) error {
	err := bc.validatePointer(model)
	if err != nil {
		return errors.Wrap(err, "unable to validate pointer")
	}

	err = bc.db.NewSelect().Model(model).Scan(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve all entities")
	}

	return nil
}

func (bc *baseClient) delete(ctx context.Context, model interface{}, primaryCol string, primaryKey interface{}) error {
	err := bc.validatePointer(model)
	if err != nil {
		return errors.Wrap(err, "unable to validate pointer")
	}
	results, err := bc.db.NewDelete().Model(model).Where("? = ?", bun.Ident(primaryCol), primaryKey).Exec(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to delete '%s':'%s'", primaryCol, primaryKey)
	}
	affected, err := results.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "unable to determine rows affected")
	}
	if affected <= 0 {
		return errors.Wrapf(storage.ErrEntityNotFound, "unable to find '%s':'%s'", primaryCol, primaryKey)
	}

	return nil
}

func (bc *baseClient) update(ctx context.Context, model interface{}, primaryCol string, primaryKey interface{}, columns ...string) error {
	err := bc.validatePointer(model)
	if err != nil {
		return errors.Wrap(err, "unable to validate pointer")
	}
	_, err = bc.db.NewUpdate().Model(model).OmitZero().Where("? = ?", bun.Ident(primaryCol), primaryKey).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to update entity")
	}
	return nil
}

func (bc *baseClient) create(ctx context.Context, model interface{}) error {
	err := bc.validatePointer(model)
	if err != nil {
		return errors.Wrap(err, "unable to validate pointer")
	}

	_, err = bc.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == 2067 {
			return errors.Wrap(storage.ErrEntityAlreadyExists, "entity already exists")
		}

		return errors.Wrap(err, "unable to create entity")
	}
	return nil
}

func (bc *baseClient) validatePointer(ptr interface{}) error {
	if ptr == nil {
		return errors.New("nil pointer")
	}
	ptrType := reflect.TypeOf(ptr)
	if ptrType.Kind() != reflect.Ptr {
		return errors.Errorf("expected pointer or slice but got %s", ptrType.Kind())
	}

	return nil
}
