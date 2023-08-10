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
	logger  Logger
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
}

func New(logger Logger) storage.Storager {
	logger.Info("Init in-memory storage")
	return &Storage{storage: nil, logger: logger}
}

func (store *Storage) Connect() error {
	store.storage = make(dataMap)
	store.logger.Info("New connect to in-memory storage")
	return nil
}

func (store *Storage) Close() error {
	store.storage = nil
	store.logger.Info("Close connection to in-memory storage")
	return nil
}

func (store *Storage) AddEvent(event storage.Event) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if !(event.ID > 0) {
		return errEventID
	}
	store.storage[event.ID] = event
	store.logger.Info("Add new event with ID: " + string(event.ID))
	return nil
}

func (store *Storage) DeleteEvent(id int64) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if !(id > 0) {
		store.logger.Error(errEventID.Error())
		return errEventID
	}

	if _, ok := store.storage[id]; ok {
		delete(store.storage, id)
	}

	store.logger.Info("Delete event with ID: " + string(id))

	return nil
}

func (store *Storage) UpdateEvent(event storage.Event) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	id := event.ID

	if _, ok := store.storage[id]; ok {
		store.storage[id] = event
		store.logger.Info("Update event with ID: " + string(id))
		return nil
	}
	store.logger.Error(errEventID.Error())
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
	store.logger.Info("Events owned by " + owner)
	return events, nil
}
