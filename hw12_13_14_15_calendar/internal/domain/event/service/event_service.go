package eventservice

import (
	"errors"

	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event"
)

type Service interface {
	Validate(eventModel *event.Event) error
}

type service struct{}

func ProvideService() Service {
	return &service{}
}

var (
	ErrRequiredTitle       = errors.New("title is required")
	ErrRequiredDescription = errors.New("title is description")
	ErrInvalidDate         = errors.New("end date end must be after start date")
)

func (s *service) Validate(eventModel *event.Event) error {
	if eventModel.Title == "" {
		return ErrRequiredTitle
	}
	if eventModel.Description == "" {
		return ErrRequiredDescription
	}
	if eventModel.EndDate.Before(eventModel.StartDate) {
		return ErrInvalidDate
	}

	return nil
}
