package console

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/app/provider"
)

var ErrIncorrectOrEmptyMigrationName = errors.New("incorrect or empty migration name")

func GetMigrateCommands() *cobra.Command {
	migrator := provider.CalendarContainer.GetMigrator()

	commandUp := &cobra.Command{
		Use:   "up",
		Short: "Apply all up",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := migrator.Up()
			if err != nil {
				return err
			}

			return nil
		},
	}

	commandDown := &cobra.Command{
		Use:   "down",
		Short: "Rollback last migration",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := migrator.Down()
			if err != nil {
				return err
			}

			return nil
		},
	}

	commandCreate := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new migration file",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return ErrIncorrectOrEmptyMigrationName
			}

			err := migrator.Create(args[0])
			if err != nil {
				return err
			}

			return nil
		},
	}

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migration manager",
	}
	migrateCmd.AddCommand(commandUp, commandDown, commandCreate)

	return migrateCmd
}
