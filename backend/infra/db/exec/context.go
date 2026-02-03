package exec

import (
	"database/sql"

	"github.com/aereal/nikki/backend/infra/db/queries"
)

type Context = queries.DBTX

var (
	_ Context = (*sql.DB)(nil)
	_ Context = (*sql.Tx)(nil)
)
