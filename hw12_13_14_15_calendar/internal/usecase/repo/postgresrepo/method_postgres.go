package postgresrepo

import (
	"context"
	"fmt"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	errapp "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/usecase/repo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/storage/postgres"
	"time"
)

// PostrgesRepo -.
type PostrgesRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *PostrgesRepo {
	return &PostrgesRepo{pg}
}

func (s *PostrgesRepo) CreateEvent(ctx context.Context, event entity.Event) error {
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

	query := `INSERT INTO events(user_id, title, descr, event_time, duration) 
			  VALUES (:userId, :title, :descr, :event_time, :duration)`

	args := map[string]interface{}{
		"user_id":    event.UserId,
		"title":      event.Title,
		"descr":      event.Desc,
		"event_time": event.EventTime,
		"duration":   event.Duration,
	}
	_, e := s.Pool.Query(ctx, query, args)
	if e != nil {
		return fmt.Errorf("failed to create: %w", e)
	}

	return nil
}

func (s *PostrgesRepo) UpdateEvent(ctx context.Context, event entity.Event) error {
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
				SET user_id = :userId, 
				    title = :title, 
				    descr = :descr, 
				    event_time = :event_time, 
				    duration = :duration,
				    notification = :notification
				WHERE id = :id`

	args := map[string]interface{}{
		"id":           event.Id,
		"user_Id":      event.UserId,
		"title":        event.Title,
		"descr":        event.Desc,
		"event_time":   event.EventTime,
		"duration":     event.Duration,
		"notification": event.Notification,
	}
	_, e := s.Pool.Exec(ctx, query, args)
	if e != nil {
		return fmt.Errorf("failed to update: %w", e)
	}

	return nil
}

// DeleteEvent Удалить (ID события);
func (s *PostrgesRepo) DeleteEvent(ctx context.Context, Id int32) error {
	_, err := s.Pool.Exec(ctx, `delete from events where id = $1`, Id)
	return err
}

// GetDailyEvents СписокСобытийНаДень (дата);
// Выводит все события, которые начинаются в заданный день
func (s *PostrgesRepo) GetDailyEvents(ctx context.Context, date time.Time) ([]entity.Event, error) {
	var events []entity.Event
	query := `SELECT id, title, descr, user_id, event_time, duration, notification
       			FROM events 
       			WHERE DATE_TRUNC('day', event_time) = :date`
	args := map[string]interface{}{
		"date": date.Day(),
	}
	rows, err := s.Pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event entity.Event
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
func (s *PostrgesRepo) GetWeeklyEvents(ctx context.Context, date time.Time) ([]entity.Event, error) {
	var events []entity.Event
	query := `SELECT id, title, descr, user_id, event_time, duration, notification
       			FROM events 
       			WHERE DATE_TRUNC('week', event_time) 
       			BETWEEN date :date AND date :date + 7 * interval '1 day'`
	args := map[string]interface{}{
		"date": date.Day(),
	}
	rows, err := s.Pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event entity.Event
		if e := rows.Scan(&event.Id, &event.Title, &event.Desc,
			&event.UserId, &event.EventTime, &event.Duration, &event.Notification); e != nil {
			return events, e
		}
		events = append(events, event)
	}
	return events, nil
}

// GetMonthlyEvents СписокСобытийНaМесяц (дата начала месяца)
// Выводит список событий за 30 дней, начиная с дня начала
func (s *PostrgesRepo) GetMonthlyEvents(ctx context.Context, date time.Time) ([]entity.Event, error) {
	var events []entity.Event
	query := `SELECT id, title, descr, user_id, event_time, duration, notification
       			FROM events 
       			WHERE DATE_TRUNC('week', event_time) 
       			BETWEEN date :date AND date :date + 30 * interval '1 day'`
	args := map[string]interface{}{
		"date": date.Day(),
	}
	rows, err := s.Pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event entity.Event
		if e := rows.Scan(&event.Id, &event.Title, &event.Desc,
			&event.UserId, &event.EventTime, &event.Duration, &event.Notification); e != nil {
			return events, e
		}
		events = append(events, event)
	}
	return events, nil
}

// IsEventTimeBusy проверка на занятость времени
func (s *PostrgesRepo) IsEventTimeBusy(event entity.Event) (bool, error) {
	//TODO: Написать проверку времени на занятость
	query := `SELECT id 
			  FROM events 
			  WHERE user_id = :userId 
			    AND event_time BETWEEN :event_time AND :end_time
			  LIMIT 1`
	args := map[string]interface{}{
		"user_id":    event.UserId,
		"event_time": event.EventTime,
		"end_time":   event.EventTime.Add(event.Duration),
	}

	rows, err := s.Pool.Query(context.Background(), query, args)
	if err != nil {
		return true, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

// EventValidate проверка ивента на валидность
func (s *PostrgesRepo) EventValidate(event entity.Event) error {
	//TODO написать ивент валидатор
	_ = event
	return nil
}
