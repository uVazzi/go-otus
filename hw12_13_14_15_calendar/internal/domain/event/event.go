package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/common"
)

type Event struct {
	UUID                  uuid.UUID
	UserUUID              uuid.UUID
	Title                 string
	Description           string
	StartDate             time.Time
	EndDate               time.Time
	DelayNotification     *int
	DelayNotificationType common.DelayType
}
