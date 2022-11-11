package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pkg/errors"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"time"
)

var (
	ErrEventTimeBusy = errors.New("event time is already busy")
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidTime   = errors.New("invalid time")
)

type Storage struct {
	db  *sql.DB
	dsn string
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	s.db, err = sql.Open("pgx", s.dsn)
	if err != nil {
		return errors.Errorf("failed to load driver: %v", err)
	}

	err = s.db.PingContext(ctx)
	if err != nil {
		return errors.Errorf("failed to connect to db: %v", err)
	}

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) Create(event storage.Event) error {
	if err := event.EventValidate(); err != nil {
		return err
	}

	check, er := s.IsEventTimeBusy(event)
	if !check {
		return ErrEventTimeBusy
	}
	if er != nil {
		return er
	}

	query := `INSERT INTO events(userid, title, descr, event_time, duration) 
			  VALUES ($1, $2, $3, $4, $5)`

	_, e := s.db.Exec(query,
		event.UserId,
		event.Title,
		event.Desc,
		event.EventTime,
		event.Duration,
	)
	if e != nil {
		return fmt.Errorf("failed to create: %w", e)
	}

	return nil
}

func (s *Storage) Update(event storage.Event) error {
	return nil
}
func (s *Storage) Delete(Id uuid.UUID) error {
	return nil
}

func (s *Storage) GetDailyEvents(date string) ([]storage.Event, error) {
	return s.eventsInTimeSpan(start, end, check), nil
}
func (s *Storage) GetWeeklyEvents(date string) ([]storage.Event, error) {
	return s.eventsInTimeSpan(start, end, check), nil
}
func (s *Storage) GetMonthlyEvents(date string) ([]storage.Event, error) {
	return s.eventsInTimeSpan(start, end, check), nil
}

func (s *Storage) IsEventTimeBusy(event storage.Event) (bool, error) {
	//TODO: Написать проверку времени на занятость
	query := `SELECT id 
			  FROM events 
			  WHERE userid = :userId 
			    AND event_time > :time  
			    AND event_time < :end_time
			  LIMIT 1`
	args := map[string]interface{}{
		"userId":   event.UserId,
		"time":     event.EventTime,
		"end_time": event.EventTime.Add(event.Duration),
	}

	rows, err := s.db.QueryContext(context.Background(), query, args)
	if err != nil {
		return true, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (s *Storage) eventsInTimeSpan(start, end, check time.Time) []storage.Event {
	var events []storage.Event

	for _, userEvents := range s.events {
		for _, event := range userEvents {
			if check.After(start) && check.Before(end) {
				events = append(events, event)
			}
		}
	}

	return events
}
