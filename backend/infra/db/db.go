package db

import (
	"database/sql"

	"github.com/XSAM/otelsql"
	"go.opentelemetry.io/otel/trace"
	_ "modernc.org/sqlite"
)

func ProvideDB(tp trace.TracerProvider, ep Endpoint) (*sql.DB, error) {
	db, err := otelsql.Open("sqlite", ep.DataSourceName(),
		otelsql.WithTracerProvider(tp),
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
