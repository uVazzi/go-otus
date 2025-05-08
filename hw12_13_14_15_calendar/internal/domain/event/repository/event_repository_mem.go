package eventrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/memory"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/txrunner"
)

type repositoryMem struct {
	storage *memory.Storage
}

func ProvideRepositoryMemory(memoryStorage *memory.Storage) Repository {
	return &repositoryMem{
		storage: memoryStorage,
	}
}

func (repo *repositoryMem) Save(_ context.Context, _ txrunner.Tx, eventModel *event.Event) error {
	repo.storage.Mutex.Lock()
	defer repo.storage.Mutex.Unlock()

	if eventModel == nil {
		return ErrEmptyEvent
	}
	if eventModel.UUID == uuid.Nil {
		eventModel.UUID = uuid.New()
	}

	repo.storage.Events[eventModel.UUID] = *eventModel
	return nil
}

func (repo *repositoryMem) Delete(_ context.Context, _ txrunner.Tx, eventUUID uuid.UUID) error {
	repo.storage.Mutex.Lock()
	defer repo.storage.Mutex.Unlock()

	_, exists := repo.storage.Events[eventUUID]
	if !exists {
		return ErrEventNotFoundByDelete
	}

	delete(repo.storage.Events, eventUUID)
	return nil
}

func (repo *repositoryMem) Get(_ context.Context, filter Filter) (*event.Event, error) {
	repo.storage.Mutex.RLock()
	defer repo.storage.Mutex.RUnlock()

	err := CheckFilter(filter)
	if err != nil {
		return nil, err
	}

	for _, eventItem := range repo.storage.Events {
		if matchFilter(eventItem, filter) {
			return &eventItem, nil
		}
	}

	return nil, ErrEventNotFoundByGet
}

func (repo *repositoryMem) List(_ context.Context, filter Filter) ([]*event.Event, error) {
	repo.storage.Mutex.RLock()
	defer repo.storage.Mutex.RUnlock()

	err := CheckFilter(filter)
	if err != nil {
		return nil, err
	}

	var result []*event.Event
	for _, eventItem := range repo.storage.Events {
		if matchFilter(eventItem, filter) {
			result = append(result, &eventItem)
		}
	}

	return result, nil
}

func (repo *repositoryMem) Exists(_ context.Context, filter Filter) (bool, error) {
	repo.storage.Mutex.RLock()
	defer repo.storage.Mutex.RUnlock()

	err := CheckFilter(filter)
	if err != nil {
		return false, err
	}

	for _, eventItem := range repo.storage.Events {
		if matchFilter(eventItem, filter) {
			return true, nil
		}
	}

	return false, nil
}

func matchFilter(event event.Event, filter Filter) bool {
	if filter.UUID != nil && event.UUID != *filter.UUID {
		return false
	}
	if filter.UserUUID != nil && event.UserUUID != *filter.UserUUID {
		return false
	}

	if filter.StartDateByBusy != nil && filter.EndDateByBusy == nil {
		if filter.StartDateByBusy.Before(event.StartDate) || filter.StartDateByBusy.After(event.EndDate) {
			return false
		}
	}
	if filter.EndDateByBusy != nil && filter.StartDateByBusy == nil {
		if filter.EndDateByBusy.Before(event.StartDate) || filter.EndDateByBusy.After(event.EndDate) {
			return false
		}
	}
	if filter.StartDateByBusy != nil && filter.EndDateByBusy != nil {
		if (filter.StartDateByBusy.Before(event.StartDate) || filter.StartDateByBusy.After(event.EndDate)) &&
			(filter.EndDateByBusy.Before(event.StartDate) || filter.EndDateByBusy.After(event.EndDate)) {
			return false
		}
	}

	return true
}
