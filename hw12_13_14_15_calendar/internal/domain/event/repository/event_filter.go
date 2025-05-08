package eventrepo

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyFilter         = errors.New("empty filter")
	ErrIncorrectDateFilter = errors.New("incorrect date filter")
)

type Filter struct {
	UUID            *uuid.UUID
	UserUUID        *uuid.UUID
	StartDateByBusy *time.Time
	EndDateByBusy   *time.Time
}

func CheckFilter(filter Filter) error {
	if filter.UUID == nil && filter.UserUUID == nil && filter.StartDateByBusy == nil && filter.EndDateByBusy == nil {
		return ErrEmptyFilter
	}

	if filter.StartDateByBusy != nil && filter.EndDateByBusy != nil &&
		filter.StartDateByBusy.After(*filter.EndDateByBusy) {
		return ErrIncorrectDateFilter
	}

	return nil
}
