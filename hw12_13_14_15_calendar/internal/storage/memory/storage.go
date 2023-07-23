package memorystorage

import (
	"errors"
	"sync"

	//"github.com/chilts/sid"

	"github.com/vpsitnik/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	storage map[string]storage.Event
	mu      sync.RWMutex //nolint:unused
}

var errEventID = errors.New("Event ID is not correct")

func New() *Storage {
	return &Storage{storage: map[string]storage.Event{}}
}

func (store *Storage) AddEvent(event storage.Event) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if event.ID == "" {
		return errEventID
	}
	store.storage[event.ID] = event

	return nil
}

func (store *Storage) DeleteEvent(id string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if id == "" {
		return errEventID
	}

	if _, ok := store.storage[id]; ok {
		delete(store.storage, id)
	}

	return nil
}

func (store *Storage) ListEventsByOwner(owner string) ([]storage.Event, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var events []storage.Event

	for _, event := range store.storage {
		if owner == event.Owner {
			events = append(events, event)
		}
	}
	return events, nil
}
