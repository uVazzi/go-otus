package txrunner

import (
	"context"
	"database/sql"

	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/config"
)

type Tx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type TxRunner interface {
	Run(ctx context.Context, useTransaction bool, fn func(ctx context.Context, tx Tx) error) error
}

type txRunner struct {
	db                 *sql.DB
	withoutTransaction bool
}

func ProvideTxRunner(config *config.Config, db *sql.DB) TxRunner {
	return &txRunner{
		db:                 db,
		withoutTransaction: config.App.UseMemoryStorage,
	}
}

func (runner *txRunner) Run(
	ctx context.Context,
	useTransaction bool,
	fn func(ctx context.Context, tx Tx) error,
) (err error) {
	if runner.withoutTransaction || !useTransaction {
		nonTx := &withoutTxAdapter{db: runner.db}
		return fn(ctx, nonTx)
	}

	tx, err := runner.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(ctx, tx)
	return err
}
