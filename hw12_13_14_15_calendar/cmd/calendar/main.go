package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/cmd/calendar/console"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/app/provider"
	internalhttp "github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/infra/http"
)

var commands = &cobra.Command{
	Use: "calendar",
	Long: `The "Calendar" service is a maximally simplified service
for storing calendar events and sending notifications`,
	Run: runMain,
}

func init() {
	var configFilePath string
	commands.PersistentFlags().StringVar(&configFilePath, "config", "configs/config.yml", "Path to configuration file")

	// Позже cobra сама сделает ParseFlags при Execute, но нам нужен config для инициализации DI провайдера
	_ = commands.ParseFlags(os.Args[1:])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	provider.ProvideContainer(ctx, configFilePath)
}

func main() {
	defer provider.CalendarContainer.CloseDB()

	commands.AddCommand(console.GetMigrateCommands())
	commands.AddCommand(console.GetVersionCommand())

	err := commands.Execute()
	if err != nil {
		provider.CalendarContainer.GetLogger().Error(context.TODO(), "Error execute: "+err.Error())
		os.Exit(1) //nolint:gocritic
	}
}

func runMain(_ *cobra.Command, _ []string) {
	logg := provider.CalendarContainer.GetLogger()
	server := internalhttp.NewServer()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error(ctx, "failed to stop http server: "+err.Error())
		}
	}()

	logg.Info(context.TODO(), "calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error(ctx, "failed to start http server: "+err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
