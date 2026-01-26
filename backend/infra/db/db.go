package db

import (
	"database/sql"

	"github.com/XSAM/otelsql"
	_ "github.com/mattn/go-sqlite3"
	"go.opentelemetry.io/otel/trace"
)

func ProvideDB(tp trace.TracerProvider, ep Endpoint) (*sql.DB, error) {
	db, err := otelsql.Open("sqlite3", ep.DataSourceName(),
		otelsql.WithTracerProvider(tp),
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
