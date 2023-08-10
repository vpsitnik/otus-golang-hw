package sqlstorage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/vpsitnik/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db     *sql.DB
	logger Logger
	ctx    context.Context
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
}

func New(dsn string, logger Logger) storage.Storager {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info("Open new SQL")
	return &Storage{db: db, logger: logger}
}

func (s *Storage) Connect() error {
	err := s.db.PingContext(s.ctx)
	if err != nil {
		s.logger.Error(err.Error())
	}

	s.logger.Info("Connect SQL DB")
	return err
}

func (s *Storage) Close() error {
	s.logger.Info("Close SQL DB")
	return s.db.Close()
}

func (s *Storage) AddEvent(event storage.Event) error {
	query := `INSERT INTO events("title", "timestamp", "duration", "description", "owner") VALUES($1, NOW(), $3, $4, $5);`
	result, err := s.db.ExecContext(s.ctx, query, event.Title, event.Timestamp, event.Duration, event.Description, event.Owner)
	if err != nil {
		s.logger.Error(err.Error())
	}
	id, _ := result.LastInsertId()
	s.logger.Info("New event with ID: " + string(id))

	return err
}

func (s *Storage) DeleteEvent(id int64) error {
	query := `DELETE FROM events WHERE id = $1;`
	result, err := s.db.ExecContext(s.ctx, query, id)
	if err != nil {
		s.logger.Error(err.Error())
	}
	rows, _ := result.RowsAffected()
	s.logger.Info("Affected rows: " + string(rows))
	return err
}

func (s *Storage) UpdateEvent(event storage.Event) error {
	query := `UPDATE events SET "title"=$2, "description"=$3, "owner"=$4 WHERE id = $1;`
	result, err := s.db.ExecContext(s.ctx, query, event.ID, event.Title, event.Description, event.Owner)
	if err != nil {
		s.logger.Error(err.Error())
	}
	rows, _ := result.RowsAffected()
	s.logger.Info("Affected rows: " + string(rows))
	return err
}

func (s *Storage) ListEventsByOwner(owner string) ([]storage.Event, error) {
	var events []storage.Event
	query := `SELECT * FROM events WHERE owner=$1;`

	result, err := s.db.ExecContext(s.ctx, query, owner)
	if err != nil {
		s.logger.Error(err.Error())
		return events, err
	}
	rows, _ := result.RowsAffected()
	s.logger.Info("Affected rows: " + string(rows))
	return events, nil
}
