package eventservice

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event"
)

func TestService(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		s := ProvideService()

		eventModel := &event.Event{
			Title:       "title",
			Description: "description",
		}
		err := s.Validate(eventModel)
		assert.NoError(t, err)

		eventModel = &event.Event{
			Description: "description",
		}
		err = s.Validate(eventModel)
		require.Truef(t, errors.Is(err, ErrRequiredTitle), "actual err - %v", err)

		eventModel = &event.Event{
			Title: "title",
		}
		err = s.Validate(eventModel)
		require.Truef(t, errors.Is(err, ErrRequiredDescription), "actual err - %v", err)

		eventModel = &event.Event{
			Title:       "title",
			Description: "description",
			StartDate:   time.Now(),
			EndDate:     time.Now().Add(-1 * time.Hour),
		}
		err = s.Validate(eventModel)
		require.Truef(t, errors.Is(err, ErrInvalidDate), "actual err - %v", err)
	})
}
