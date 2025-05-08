package txrunner

import (
	"context"
	"database/sql"
)

type withoutTxAdapter struct {
	db *sql.DB
}

func (adapter *withoutTxAdapter) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return adapter.db.ExecContext(ctx, query, args...)
}

func (adapter *withoutTxAdapter) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return adapter.db.QueryContext(ctx, query, args...)
}

func (adapter *withoutTxAdapter) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return adapter.db.QueryRowContext(ctx, query, args...)
}
