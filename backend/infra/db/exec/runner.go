package exec

import (
	"context"
	"database/sql"
)

func ProvideRunner(db *sql.DB) *Runner {
	return &Runner{db: db}
}

type Runner struct {
	db *sql.DB
}

var _ Context = (*Runner)(nil)

func (r *Runner) getContext(ctx context.Context) Context {
	if tx, ok := getTxFromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *Runner) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return r.getContext(ctx).ExecContext(ctx, query, args...)
}

func (r *Runner) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return r.getContext(ctx).PrepareContext(ctx, query)
}

func (r *Runner) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return r.getContext(ctx).QueryContext(ctx, query, args...)
}

func (r *Runner) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return r.getContext(ctx).QueryRowContext(ctx, query, args...)
}

type keyTx struct{}

func getTxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(keyTx{}).(*sql.Tx)
	return tx, ok
}

func contextWithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, keyTx{}, tx)
}
