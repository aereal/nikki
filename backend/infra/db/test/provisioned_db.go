package test

import (
	"context"
	"database/sql"

	"github.com/aereal/nikki/backend/infra/db"
)

func provideProvisionedDB(ctx context.Context, sqldb *sql.DB) (*provisionedDB, error) {
	if _, err := sqldb.ExecContext(ctx, db.Schema()); err != nil {
		return nil, err
	}
	return &provisionedDB{db: sqldb}, nil
}

type provisionedDB struct {
	db *sql.DB
}
