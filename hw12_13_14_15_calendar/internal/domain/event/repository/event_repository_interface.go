package eventrepo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/txrunner"
)

var (
	ErrEmptyEvent            = errors.New("empty event")
	ErrEventNotFoundByGet    = errors.New("event not found by get")
	ErrEventNotFoundByDelete = errors.New("event not found by delete")
)

type Repository interface {
	Save(ctx context.Context, tx txrunner.Tx, eventModel *event.Event) error
	Delete(ctx context.Context, tx txrunner.Tx, eventUUID uuid.UUID) error
	Get(ctx context.Context, filter Filter) (*event.Event, error)
	List(ctx context.Context, filter Filter) ([]*event.Event, error)
	Exists(ctx context.Context, filter Filter) (bool, error)
}
