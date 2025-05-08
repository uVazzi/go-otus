package eventusecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/common"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event"
	eventrepo "github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event/repository"
	eventservice "github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event/service"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/logger"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/txrunner"
	eventdto "github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/usecase/event/dto"
)

type eventUsecase struct {
	tx        txrunner.TxRunner
	eventRepo eventrepo.Repository
	service   eventservice.Service
	logger    *logger.Logger
}

type EventUsecase interface {
	Create(ctx context.Context, dto eventdto.EventCreateDTO) error
}

func ProvideUsecase(
	tx txrunner.TxRunner,
	logger *logger.Logger,
	repo eventrepo.Repository,
	service eventservice.Service,
) EventUsecase {
	return &eventUsecase{
		tx:        tx,
		eventRepo: repo,
		service:   service,
		logger:    logger,
	}
}

var (
	ErrDateBusy          = errors.New("date busy")
	ErrUserRequired      = errors.New("user required")
	ErrIncorrectUserUUID = errors.New("incorrect user uuid")
	ErrUserNotFound      = errors.New("user not found")
)

func (uc *eventUsecase) Create(ctx context.Context, dto eventdto.EventCreateDTO) error {
	eventModel, err := getEventForCreate(dto)
	if err != nil {
		return err
	}

	err = uc.service.Validate(eventModel)
	if err != nil {
		return err
	}
	err = uc.checkUser(eventModel.UserUUID)
	if err != nil {
		return err
	}
	err = uc.checkBusyDate(ctx, eventModel.UserUUID, eventModel.StartDate, eventModel.EndDate)
	if err != nil {
		return err
	}

	err = uc.tx.Run(ctx, true, func(ctx context.Context, tx txrunner.Tx) error {
		return uc.eventRepo.Save(ctx, tx, eventModel)
	})
	if err != nil {
		return err
	}

	return nil
}

func (uc *eventUsecase) checkUser(userUUID uuid.UUID) error {
	_ = uc

	UUID, _ := uuid.Parse("8f3d77cc-1234-4567-890a-abcdefabcdef")
	if UUID != userUUID {
		return ErrUserNotFound
	}

	return nil
}

func (uc *eventUsecase) checkBusyDate(ctx context.Context, userUUID uuid.UUID, startDate, endDate time.Time) error {
	searchFilter := eventrepo.Filter{
		UserUUID:        &userUUID,
		StartDateByBusy: &startDate,
		EndDateByBusy:   &endDate,
	}

	exists, err := uc.eventRepo.Exists(ctx, searchFilter)
	if err != nil {
		return err
	}

	if exists {
		return ErrDateBusy
	}

	return nil
}

func getEventForCreate(dto eventdto.EventCreateDTO) (*event.Event, error) {
	if dto.UserUUID == "" {
		return nil, ErrUserRequired
	}
	userUUID, err := uuid.Parse(dto.UserUUID)
	if err != nil {
		return nil, ErrIncorrectUserUUID
	}

	delayType := common.DelayType("")
	if dto.DelayNotification != nil {
		delayType, err = common.GetDelayType(dto.DelayNotificationType)
		if err != nil {
			return nil, err
		}
	}

	eventModel := &event.Event{
		UserUUID:              userUUID,
		Title:                 dto.Title,
		Description:           dto.Description,
		StartDate:             dto.StartDate,
		EndDate:               dto.EndDate,
		DelayNotification:     dto.DelayNotification,
		DelayNotificationType: delayType,
	}

	return eventModel, nil
}
