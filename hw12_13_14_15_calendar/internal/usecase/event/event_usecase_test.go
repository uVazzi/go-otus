package eventusecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/common"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event"
	eventrepo "github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event/repository"
	eventservice "github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event/service"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/config"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/logger"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/memory"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/txrunner"
	eventdto "github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/usecase/event/dto"
)

func TestUsecaseCreate(t *testing.T) {
	conf := &config.Config{
		App: config.AppConf{
			UseMemoryStorage: true,
		},
		Logger: config.LoggerConf{
			LogLevel:        "DEBUG",
			DisableErrorLog: false,
			DisableWarnLog:  false,
			DisableInfoLog:  false,
			DisableDebugLog: false,
		},
	}

	logg := logger.ProvideLogger(conf)
	txRunner := txrunner.ProvideTxRunner(conf, &sql.DB{})

	memoryStorage := memory.ProvideMemory()
	repo := eventrepo.ProvideRepositoryMemory(memoryStorage)
	servece := eventservice.ProvideService()

	us := ProvideUsecase(txRunner, logg, repo, servece)

	t.Run("success", func(t *testing.T) {
		eventDTO := eventdto.EventCreateDTO{
			UserUUID:    "8f3d77cc-1234-4567-890a-abcdefabcdef",
			Title:       "title",
			Description: "description",
			StartDate:   time.Now().Add(time.Hour),
			EndDate:     time.Now().Add(2 * time.Hour),
		}
		err := us.Create(context.TODO(), eventDTO)
		assert.NoError(t, err)
	})

	t.Run("check getEventForCreate", func(t *testing.T) {
		eventDTO := eventdto.EventCreateDTO{
			UserUUID: "",
		}
		err := us.Create(context.TODO(), eventDTO)
		require.Truef(t, errors.Is(err, ErrUserRequired), "actual err - %v", err)

		eventDTO = eventdto.EventCreateDTO{
			UserUUID: "INCORRECT_UUID",
		}
		err = us.Create(context.TODO(), eventDTO)
		require.Truef(t, errors.Is(err, ErrIncorrectUserUUID), "actual err - %v", err)

		delayNotification := 2
		eventDTO = eventdto.EventCreateDTO{
			UserUUID:              "8f3d77cc-1234-4567-890a-abcdefabcdef",
			DelayNotification:     &delayNotification,
			DelayNotificationType: "year",
		}
		err = us.Create(context.TODO(), eventDTO)
		require.Truef(t, errors.Is(err, common.ErrDelayType), "actual err - %v", err)
	})

	t.Run("check checkUser", func(t *testing.T) {
		eventDTO := eventdto.EventCreateDTO{
			UserUUID:    "6f3d77cc-1234-4567-890a-abcdefabcdef",
			Title:       "title",
			Description: "description",
			StartDate:   time.Now().Add(time.Hour),
			EndDate:     time.Now().Add(2 * time.Hour),
		}
		err := us.Create(context.TODO(), eventDTO)
		require.Truef(t, errors.Is(err, ErrUserNotFound), "actual err - %v", err)
	})

	t.Run("check checkUser", func(t *testing.T) {
		userUUID, _ := uuid.Parse("8f3d77cc-1234-4567-890a-abcdefabcdef")
		eventModel := &event.Event{
			UserUUID:    userUUID,
			Title:       "title",
			Description: "description",
			StartDate:   time.Now().Add(time.Hour),
			EndDate:     time.Now().Add(2 * time.Hour),
		}
		err := repo.Save(context.TODO(), nil, eventModel)
		assert.NoError(t, err)

		eventDTO := eventdto.EventCreateDTO{
			UserUUID:    "8f3d77cc-1234-4567-890a-abcdefabcdef",
			Title:       "title",
			Description: "description",
			StartDate:   time.Now().Add(time.Hour),
			EndDate:     time.Now().Add(2 * time.Hour),
		}
		err = us.Create(context.TODO(), eventDTO)
		require.Truef(t, errors.Is(err, ErrDateBusy), "actual err - %v", err)
	})
}
