package eventrepo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/memory"
)

func TestRepository(t *testing.T) {
	storage := memory.ProvideMemory()
	repo := ProvideRepositoryMemory(storage)

	userUUIDBy12 := uuid.New()
	eventModel1 := &event.Event{
		UUID:        uuid.New(),
		UserUUID:    userUUIDBy12,
		Title:       "Event 1",
		Description: "Description event 1",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour),
	}
	eventModel2 := &event.Event{
		UUID:        uuid.New(),
		UserUUID:    userUUIDBy12,
		Title:       "Event 2",
		Description: "Description event 2",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(2 * time.Hour),
	}
	eventModel3 := &event.Event{
		UUID:        uuid.New(),
		UserUUID:    uuid.New(),
		Title:       "Event 3",
		Description: "Description event 3",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(3 * time.Hour),
	}

	t.Run("save empty event", func(t *testing.T) {
		err := repo.Save(context.TODO(), nil, nil)
		assert.ErrorIs(t, err, ErrEmptyEvent)
	})

	t.Run("not found get", func(t *testing.T) {
		testUUID := uuid.New()
		_, err := repo.Get(context.TODO(), Filter{UUID: &testUUID})

		require.Truef(t, errors.Is(err, ErrEventNotFoundByGet), "actual err - %v", err)
	})

	t.Run("not found delete", func(t *testing.T) {
		err := repo.Delete(context.TODO(), nil, uuid.New())
		assert.ErrorIs(t, err, ErrEventNotFoundByDelete)
	})

	t.Run("save and get by uuid", func(t *testing.T) {
		err := repo.Save(context.TODO(), nil, eventModel1)
		assert.NoError(t, err)

		foundEventModel, err := repo.Get(context.TODO(), Filter{UUID: &eventModel1.UUID})
		assert.NoError(t, err)
		assert.Equal(t, eventModel1, foundEventModel)

		eventModel1.Title = "Check change model no save"
		foundEventModel, err = repo.Get(context.TODO(), Filter{UUID: &eventModel1.UUID})
		assert.NoError(t, err)
		assert.Equal(t, "Event 1", foundEventModel.Title)

		foundEventModel.Title = "Check change found model no save"
		foundEventModel, err = repo.Get(context.TODO(), Filter{UUID: &eventModel1.UUID})
		assert.NoError(t, err)
		assert.Equal(t, "Event 1", foundEventModel.Title)
	})

	t.Run("list", func(t *testing.T) {
		err := repo.Save(context.TODO(), nil, eventModel1)
		assert.NoError(t, err)
		err = repo.Save(context.TODO(), nil, eventModel2)
		assert.NoError(t, err)
		err = repo.Save(context.TODO(), nil, eventModel3)
		assert.NoError(t, err)

		foundList, err := repo.List(context.TODO(), Filter{UserUUID: &userUUIDBy12})

		assert.NoError(t, err)
		assert.Len(t, foundList, 2)
	})

	t.Run("list", func(t *testing.T) {
		err := repo.Save(context.TODO(), nil, eventModel3)
		assert.NoError(t, err)

		_, err = repo.Get(context.TODO(), Filter{UUID: &eventModel3.UUID})
		assert.NoError(t, err)

		err = repo.Delete(context.TODO(), nil, eventModel3.UUID)
		assert.NoError(t, err)

		_, err = repo.Get(context.TODO(), Filter{UUID: &eventModel3.UUID})
		require.Truef(t, errors.Is(err, ErrEventNotFoundByGet), "actual err - %v", err)
	})
}

func TestFilter(t *testing.T) {
	storage := memory.ProvideMemory()
	repo := ProvideRepositoryMemory(storage)

	t.Run("check filter", func(t *testing.T) {
		_, err := repo.Get(context.TODO(), Filter{})
		require.Truef(t, errors.Is(err, ErrEmptyFilter), "actual err - %v", err)
		_, err = repo.List(context.TODO(), Filter{})
		require.Truef(t, errors.Is(err, ErrEmptyFilter), "actual err - %v", err)
		_, err = repo.Exists(context.TODO(), Filter{})
		require.Truef(t, errors.Is(err, ErrEmptyFilter), "actual err - %v", err)

		filterDate1 := time.Now().Add(4 * time.Hour)
		filterDate2 := time.Now().Add(2 * time.Hour)
		_, err = repo.Exists(context.TODO(), Filter{StartDateByBusy: &filterDate1, EndDateByBusy: &filterDate2})
		require.Truef(t, errors.Is(err, ErrIncorrectDateFilter), "actual err - %v", err)
	})
}

func TestBusyDate(t *testing.T) {
	storage := memory.ProvideMemory()
	repo := ProvideRepositoryMemory(storage)

	eventModel := &event.Event{
		UUID:        uuid.New(),
		UserUUID:    uuid.New(),
		Title:       "Event 3",
		Description: "Description event 3",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(3 * time.Hour),
	}

	t.Run("exists busy date", func(t *testing.T) {
		err := repo.Save(context.TODO(), nil, eventModel)
		assert.NoError(t, err)

		filterDate := time.Now().Add(-1 * time.Hour)
		exists, err := repo.Exists(context.TODO(), Filter{StartDateByBusy: &filterDate})
		assert.NoError(t, err)
		assert.False(t, exists)

		filterDate = time.Now().Add(5 * time.Hour)
		exists, err = repo.Exists(context.TODO(), Filter{StartDateByBusy: &filterDate})
		assert.NoError(t, err)
		assert.False(t, exists)

		filterDate = time.Now().Add(time.Hour)
		exists, err = repo.Exists(context.TODO(), Filter{StartDateByBusy: &filterDate})
		assert.NoError(t, err)
		assert.True(t, exists)

		filterDate = time.Now().Add(-1 * time.Hour)
		exists, err = repo.Exists(context.TODO(), Filter{EndDateByBusy: &filterDate})
		assert.NoError(t, err)
		assert.False(t, exists)

		filterDate = time.Now().Add(5 * time.Hour)
		exists, err = repo.Exists(context.TODO(), Filter{EndDateByBusy: &filterDate})
		assert.NoError(t, err)
		assert.False(t, exists)

		filterDate = time.Now().Add(time.Hour)
		exists, err = repo.Exists(context.TODO(), Filter{EndDateByBusy: &filterDate})
		assert.NoError(t, err)
		assert.True(t, exists)

		filterDate1 := time.Now().Add(-2 * time.Hour)
		filterDate2 := time.Now().Add(-1 * time.Hour)
		exists, err = repo.Exists(context.TODO(), Filter{StartDateByBusy: &filterDate1, EndDateByBusy: &filterDate2})
		assert.NoError(t, err)
		assert.False(t, exists)

		filterDate1 = time.Now().Add(4 * time.Hour)
		filterDate2 = time.Now().Add(5 * time.Hour)
		exists, err = repo.Exists(context.TODO(), Filter{StartDateByBusy: &filterDate1, EndDateByBusy: &filterDate2})
		assert.NoError(t, err)
		assert.False(t, exists)

		filterDate1 = time.Now().Add(-1 * time.Hour)
		filterDate2 = time.Now().Add(1 * time.Hour)
		exists, err = repo.Exists(context.TODO(), Filter{StartDateByBusy: &filterDate1, EndDateByBusy: &filterDate2})
		assert.NoError(t, err)
		assert.True(t, exists)

		filterDate1 = time.Now().Add(2 * time.Hour)
		filterDate2 = time.Now().Add(4 * time.Hour)
		exists, err = repo.Exists(context.TODO(), Filter{StartDateByBusy: &filterDate1, EndDateByBusy: &filterDate2})
		assert.NoError(t, err)
		assert.True(t, exists)

		filterDate1 = time.Now().Add(1 * time.Hour)
		filterDate2 = time.Now().Add(2 * time.Hour)
		exists, err = repo.Exists(context.TODO(), Filter{StartDateByBusy: &filterDate1, EndDateByBusy: &filterDate2})
		assert.NoError(t, err)
		assert.True(t, exists)
	})
}
