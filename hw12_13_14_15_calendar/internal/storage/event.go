package storage

type Event struct {
	ID          string
	Title       string
	Timestamp   int64
	Duration    int64
	Description string
	Owner       string
}
