package db

import (
	"database/sql"

	"github.com/XSAM/otelsql"
	"go.opentelemetry.io/otel/trace"
	_ "modernc.org/sqlite"
)

func ProvideDB(tp trace.TracerProvider, ep Endpoint) (*sql.DB, error) {
	dsn, err := ep.DataSourceName()
	if err != nil {
		return nil, err
	}
	db, err := otelsql.Open("sqlite", dsn,
		otelsql.WithTracerProvider(tp),
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
