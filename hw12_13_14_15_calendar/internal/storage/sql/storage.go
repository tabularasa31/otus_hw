package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pkg/errors"
	errapp "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/domain/errors"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/domain/models"
	"time"
)

type Storage struct {
	db  *sql.DB
	dsn string
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
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

func (s *Storage) Create(ctx context.Context, event models.Event) error {
	if err := s.EventValidate(event); err != nil {
		return err
	}

	check, er := s.IsEventTimeBusy(event)
	if !check {
		return errapp.ErrEventTimeBusy
	}
	if er != nil {
		return er
	}

	query := `INSERT INTO events(userid, title, descr, event_time, duration) 
			  VALUES (:userId, :title, :descr, :event_time, :duration)`

	args := map[string]interface{}{
		"userId":     event.UserId,
		"title":      event.Title,
		"descr":      event.Desc,
		"event_time": event.EventTime,
		"duration":   event.Duration,
	}
	_, e := s.db.QueryContext(ctx, query, args)
	if e != nil {
		return fmt.Errorf("failed to create: %w", e)
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, event models.Event) error {
	if err := s.EventValidate(event); err != nil {
		return err
	}

	check, er := s.IsEventTimeBusy(event)
	if !check {
		return errapp.ErrEventTimeBusy
	}
	if er != nil {
		return er
	}

	query := `UPDATE events 
				SET userid = :userId, 
				    title = :title, 
				    descr = :descr, 
				    event_time = :event_time, 
				    duration = :duration,
				    notification = :notification
				WHERE id = :id`

	args := map[string]interface{}{
		"id":           event.Id,
		"userId":       event.UserId,
		"title":        event.Title,
		"descr":        event.Desc,
		"event_time":   event.EventTime,
		"duration":     event.Duration,
		"notification": event.Notification,
	}
	_, e := s.db.ExecContext(ctx, query, args)
	if e != nil {
		return fmt.Errorf("failed to update: %w", e)
	}

	return nil
}

// Delete Удалить (ID события);
func (s *Storage) Delete(Id int32) error {
	_, err := s.db.Exec(`delete from events where id = $1`, Id)
	return err
}

// GetDailyEvents СписокСобытийНаДень (дата);
// Выводит все события, которые начинаются в заданный день
func (s *Storage) GetDailyEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	var events []models.Event
	query := `SELECT id, title, descr, user_id, event_time, duration, notification
       			FROM events 
       			WHERE DATE_TRUNC('day', event_time) = :date`
	args := map[string]interface{}{
		"date": date.Day(),
	}
	rows, err := s.db.QueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			fmt.Errorf("failed to close row: %w", closeErr)
		}
	}()

	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.Id, &event.Title, &event.Desc,
			&event.UserId, &event.EventTime, &event.Duration, &event.Notification); err != nil {
			return events, err
		}
		events = append(events, event)
	}

	return events, nil
}

// GetWeeklyEvents СписокСобытийНаНеделю (дата начала недели);
// Выводит список событий за 7 дней, начиная с дня начала
func (s *Storage) GetWeeklyEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	var events []models.Event
	query := `SELECT id, title, descr, user_id, event_time, duration, notification
       			FROM events 
       			WHERE DATE_TRUNC('week', event_time) 
       			BETWEEN date :date AND date :date + 7 * interval '1 day'`
	args := map[string]interface{}{
		"date": date.Day(),
	}
	rows, err := s.db.QueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			fmt.Errorf("failed to close row: %w", closeErr)
		}
	}()
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.Id, &event.Title, &event.Desc,
			&event.UserId, &event.EventTime, &event.Duration, &event.Notification); err != nil {
			return events, err
		}
		events = append(events, event)
	}
	return events, nil
}

// GetMonthlyEvents СписокСобытийНaМесяц (дата начала месяца)
// Выводит список событий за 30 дней, начиная с дня начала
func (s *Storage) GetMonthlyEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	var events []models.Event
	query := `SELECT id, title, descr, user_id, event_time, duration, notification
       			FROM events 
       			WHERE DATE_TRUNC('week', event_time) 
       			BETWEEN date :date AND date :date + 30 * interval '1 day'`
	args := map[string]interface{}{
		"date": date.Day(),
	}
	rows, err := s.db.QueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			fmt.Errorf("failed to close row: %w", closeErr)
		}
	}()

	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.Id, &event.Title, &event.Desc,
			&event.UserId, &event.EventTime, &event.Duration, &event.Notification); err != nil {
			return events, err
		}
		events = append(events, event)
	}
	return events, nil
}

// IsEventTimeBusy проверка на занятость времени
func (s *Storage) IsEventTimeBusy(event models.Event) (bool, error) {
	//TODO: Написать проверку времени на занятость
	query := `SELECT id 
			  FROM events 
			  WHERE userid = :userId 
			    AND event_time BETWEEN :event_time AND :end_time
			  LIMIT 1`
	args := map[string]interface{}{
		"userId":     event.UserId,
		"event_time": event.EventTime,
		"end_time":   event.EventTime.Add(event.Duration),
	}

	rows, err := s.db.QueryContext(context.Background(), query, args)
	if err != nil {
		return true, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			fmt.Errorf("failed to close row: %w", closeErr)
		}
	}()

	return rows.Next(), nil
}

// EventValidate проверка ивента на валидность
func (s *Storage) EventValidate(event models.Event) error {
	//TODO написать ивент валидатор
	_ = event
	return nil
}
