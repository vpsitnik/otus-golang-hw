package memorystorage

import (
	"errors"
	"sync"

	"github.com/vpsitnik/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

var errEventID = errors.New("Event ID is not correct")

type dataMap map[int64]storage.Event

type Storage struct {
	storage map[int64]storage.Event
	mu      sync.RWMutex //nolint:unused
}

func New() storage.Storager {
	return &Storage{storage: nil}
}

func (store *Storage) Connect() error {
	store.storage = make(dataMap)
	return nil
}

func (store *Storage) Close() error {
	store.storage = nil
	return nil
}

func (store *Storage) AddEvent(event storage.Event) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if !(event.ID > 0) {
		return errEventID
	}
	store.storage[event.ID] = event

	return nil
}

func (store *Storage) DeleteEvent(id int64) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if !(id > 0) {
		return errEventID
	}

	if _, ok := store.storage[id]; ok {
		delete(store.storage, id)
	}

	return nil
}

func (store *Storage) UpdateEvent(event storage.Event) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	id := event.ID

	if _, ok := store.storage[id]; ok {
		store.storage[id] = event
	}

	return errEventID
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
