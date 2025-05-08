package provider

import (
	"context"
	"database/sql"
	"os"

	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/config"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/database"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/dbmigrator"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/logger"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/memory"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/txrunner"
)

var CalendarContainer *appContainer

type appContainer struct {
	conf          *config.Config
	logg          *logger.Logger
	db            *sql.DB
	migrator      dbmigrator.Migrator
	txRunner      txrunner.TxRunner
	memoryStorage *memory.Storage
	dependency    *Dependencies
}

func ProvideContainer(ctx context.Context, configFilePath string) {
	// Temp provide for init errors
	logg := logger.ProvideLogger(&config.Config{
		App: config.AppConf{
			Stage: "init",
		},
		Logger: config.LoggerConf{
			LogLevel: "ERROR",
		},
	})

	conf, err := config.ProvideConfig(configFilePath)
	if err != nil {
		logg.Error(context.TODO(), "Fatal error on read config: "+err.Error())
		os.Exit(1)
	}

	db, err := database.ProvideDB(ctx, conf)
	if err != nil {
		logg.Error(context.TODO(), "Fatal error on open db: "+err.Error())
		os.Exit(1)
	}

	migrator := dbmigrator.ProvideGoose(conf, db)
	txRunner := txrunner.ProvideTxRunner(conf, db)
	memoryStorage := memory.ProvideMemory()
	logg = logger.ProvideLogger(conf)
	dependency := ProvideDependency(conf, logg, memoryStorage, db, txRunner)

	CalendarContainer = &appContainer{
		conf:          conf,
		logg:          logg,
		db:            db,
		migrator:      migrator,
		txRunner:      txRunner,
		memoryStorage: memoryStorage,
		dependency:    dependency,
	}
}

func (container *appContainer) CloseDB() {
	if container.db.Stats().Idle != 0 || container.db.Stats().InUse != 0 {
		err := container.db.Close()
		if err != nil {
			container.logg.Error(context.TODO(), "Fatal error on close db: "+err.Error())
			os.Exit(1)
		}
	}
}

func (container *appContainer) GetMigrator() dbmigrator.Migrator {
	return container.migrator
}

func (container *appContainer) GetLogger() *logger.Logger {
	return container.logg
}

func (container *appContainer) GetConfig() *config.Config {
	return container.conf
}

func (container *appContainer) GetDependencies() *Dependencies {
	return container.dependency
}
