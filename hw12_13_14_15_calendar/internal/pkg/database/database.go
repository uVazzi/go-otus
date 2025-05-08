package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // DB driver
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/config"
)

func ProvideDB(ctx context.Context, config *config.Config) (*sql.DB, error) {
	if config.App.UseMemoryStorage {
		return &sql.DB{}, nil
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.Username,
		config.Database.Password,
		config.Database.Database,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Hour)

	return db, nil
}
