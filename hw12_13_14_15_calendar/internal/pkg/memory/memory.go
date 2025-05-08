package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event"
)

type Storage struct {
	Mutex  sync.RWMutex
	Events map[uuid.UUID]event.Event
}

func ProvideMemory() *Storage {
	return &Storage{
		Mutex:  sync.RWMutex{},
		Events: make(map[uuid.UUID]event.Event),
	}
}
