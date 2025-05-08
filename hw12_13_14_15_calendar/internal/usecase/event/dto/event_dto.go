package eventdto

import (
	"time"
)

type EventCreateDTO struct {
	UserUUID              string
	Title                 string
	Description           string
	StartDate             time.Time
	EndDate               time.Time
	DelayNotification     *int
	DelayNotificationType string
}
