package provider

import (
	"database/sql"

	eventrepo "github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event/repository"
	eventservice "github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event/service"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/config"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/logger"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/memory"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/txrunner"
	eventusecase "github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/usecase/event"
)

type Dependencies struct {
	Repository repository
	Service    service
	Usecase    usecase
}

type repository struct {
	EventRepository eventrepo.Repository
}

type service struct {
	EventService eventservice.Service
}

type usecase struct {
	EventUsecase eventusecase.EventUsecase
}

func ProvideDependency(
	config *config.Config,
	logger *logger.Logger,
	memoryStorage *memory.Storage,
	db *sql.DB,
	txRunner txrunner.TxRunner,
) *Dependencies {
	// Repository
	var eventRepository eventrepo.Repository
	if !config.App.UseMemoryStorage {
		eventRepository = eventrepo.ProvideRepositoryDB(db)
	} else {
		eventRepository = eventrepo.ProvideRepositoryMemory(memoryStorage)
	}

	repositories := repository{
		EventRepository: eventRepository,
	}

	// Service
	eventService := eventservice.ProvideService()

	services := service{
		EventService: eventService,
	}

	// Usecase
	eventUsecase := eventusecase.ProvideUsecase(txRunner, logger, repositories.EventRepository, services.EventService)

	usecases := usecase{
		EventUsecase: eventUsecase,
	}

	return &Dependencies{
		Repository: repositories,
		Service:    services,
		Usecase:    usecases,
	}
}
