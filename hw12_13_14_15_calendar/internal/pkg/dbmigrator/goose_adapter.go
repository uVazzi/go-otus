package dbmigrator

import (
	"database/sql"
	"errors"

	"github.com/pressly/goose/v3"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/config"
)

var ErrNotDB = errors.New("applying migrations is not possible (in-memory storage)")

type Migrator interface {
	Create(name string) error
	Up() error
	Down() error
}

type gooseAdapter struct {
	db               *sql.DB
	migrationDir     string
	useMemoryStorage bool
}

func ProvideGoose(config *config.Config, db *sql.DB) Migrator {
	return &gooseAdapter{
		db:               db,
		migrationDir:     "migrations",
		useMemoryStorage: config.App.UseMemoryStorage,
	}
}

func (adapter *gooseAdapter) Create(name string) error {
	err := goose.Create(nil, adapter.migrationDir, name, "sql")
	if err != nil {
		return err
	}
	return nil
}

func (adapter *gooseAdapter) Up() error {
	if adapter.useMemoryStorage {
		return ErrNotDB
	}

	err := goose.Up(adapter.db, adapter.migrationDir)
	if err != nil {
		return err
	}
	return nil
}

func (adapter *gooseAdapter) Down() error {
	if adapter.useMemoryStorage {
		return ErrNotDB
	}

	err := goose.Down(adapter.db, adapter.migrationDir)
	if err != nil {
		return err
	}
	return nil
}
