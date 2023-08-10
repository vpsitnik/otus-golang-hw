package storage

type Event struct {
	ID          int64
	Title       string
	Timestamp   int64
	Duration    int64
	Description string
	Owner       string
}

type Storager interface {
	Connect() error
	Close() error
	AddEvent(Event) error
	DeleteEvent(int64) error
	UpdateEvent(Event) error
	ListEventsByOwner(string) ([]Event, error)
}
