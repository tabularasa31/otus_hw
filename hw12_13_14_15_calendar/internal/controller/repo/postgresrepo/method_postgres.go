package postgresrepo

import (
	"context"
	"fmt"
	errapp "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/storage/postgres"
)

// EventRepo -.
type EventRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *EventRepo {
	return &EventRepo{pg}
}

// CreateEvent Создать (событие).
func (r *EventRepo) CreateEvent(ctx context.Context, eventDB *entity.EventDB) (*entity.Event, error) {
	check, er := r.isEventTimeBusy(*eventDB)
	if !check {
		return nil, errapp.ErrEventTimeBusy
	}
	if er != nil {
		return nil, er
	}

	query := `INSERT INTO events(user_id, title, descr, event_time, duration) 
			  VALUES (:userId, :title, :descr, :event_time, :duration)
			  RETURNING id`

	args := map[string]interface{}{
		"user_id":    eventDB.UserID,
		"title":      eventDB.Title,
		"descr":      eventDB.Desc,
		"event_time": eventDB.EventTime,
		"duration":   eventDB.Duration,
	}
	resID, e := r.Postgres.Pool.Exec(ctx, query, args)
	if e != nil {
		return nil, fmt.Errorf("failed to create: %w", e)
	}

	res, err := r.result(ctx, resID.String())
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateEvent Обновить (событие).
func (r *EventRepo) UpdateEvent(ctx context.Context, eventDB *entity.EventDB) (*entity.Event, error) {
	check, er := r.isEventTimeBusy(*eventDB)
	if !check {
		return nil, errapp.ErrEventTimeBusy
	}
	if er != nil {
		return nil, er
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
		"id":           eventDB.ID,
		"user_Id":      eventDB.UserID,
		"title":        eventDB.Title,
		"descr":        eventDB.Desc,
		"event_time":   eventDB.EventTime,
		"duration":     eventDB.Duration,
		"notification": eventDB.Notification,
	}

	resID, e := r.Postgres.Pool.Exec(ctx, query, args)
	if e != nil {
		return nil, fmt.Errorf("failed to update: %w", e)
	}

	res, err := r.result(ctx, resID.String())
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteEvent Удалить (ID события).
func (r *EventRepo) DeleteEvent(ctx context.Context, id int32) error {
	_, err := r.Postgres.Pool.Exec(ctx, `delete from events where id = $1`, id)
	return err
}

// GetDailyEvents СписокСобытийНаДень (дата).
// Выводит все события, которые начинаются в заданный день.
func (r *EventRepo) GetDailyEvents(ctx context.Context, eventDB *entity.EventDB) ([]entity.Event, error) {
	var events []entity.Event

	query := `SELECT id, title, descr, user_id, event_time, duration, notification
       			FROM events 
       			WHERE DATE_TRUNC('day', event_time) = :date`
	args := map[string]interface{}{
		"date": eventDB.EventTime.Day(),
	}

	rows, err := r.Postgres.Pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event entity.EventDB
		if er := rows.Scan(&event.ID, &event.Title, &event.Desc,
			&event.UserID, &event.EventTime, &event.Duration, &event.Notification); er != nil {
			return events, er
		}
		events = append(events, *event.Dto())
	}

	return events, nil
}

// GetWeeklyEvents СписокСобытийНаНеделю (дата начала недели).
// Выводит список событий за 7 дней, начиная с дня начала.
func (r *EventRepo) GetWeeklyEvents(ctx context.Context, eventDB *entity.EventDB) ([]entity.Event, error) {
	var events []entity.Event
	query := `SELECT id, title, descr, user_id, event_time, duration, notification
       			FROM events 
       			WHERE DATE_TRUNC('week', event_time) 
       			BETWEEN date :date AND date :date + 7 * interval '1 day'`
	args := map[string]interface{}{
		"date": eventDB.EventTime.Day(),
	}
	rows, err := r.Postgres.Pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event entity.EventDB
		if e := rows.Scan(&event.ID, &event.Title, &event.Desc,
			&event.UserID, &event.EventTime, &event.Duration, &event.Notification); e != nil {
			return events, e
		}
		events = append(events, *event.Dto())
	}
	return events, nil
}

// GetMonthlyEvents СписокСобытийНaМесяц (дата начала месяца).
// Выводит список событий за 30 дней, начиная с дня начала.
func (r *EventRepo) GetMonthlyEvents(ctx context.Context, eventDB *entity.EventDB) ([]entity.Event, error) {
	var events []entity.Event
	query := `SELECT id, title, descr, user_id, event_time, duration, notification
       			FROM events 
       			WHERE DATE_TRUNC('week', event_time) 
       			BETWEEN date :date AND date :date + 30 * interval '1 day'`
	args := map[string]interface{}{
		"date": eventDB.EventTime.Day(),
	}
	rows, err := r.Postgres.Pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event entity.EventDB
		if e := rows.Scan(&event.ID, &event.Title, &event.Desc,
			&event.UserID, &event.EventTime, &event.Duration, &event.Notification); e != nil {
			return events, e
		}
		events = append(events, *event.Dto())
	}
	return events, nil
}

// isEventTimeBusy проверка на занятость времени.
func (r *EventRepo) isEventTimeBusy(event entity.EventDB) (bool, error) {
	query := `SELECT id 
			  FROM events 
			  WHERE user_id = :userId 
			    AND event_time BETWEEN :event_time AND :end_time
			  LIMIT 1`
	args := map[string]interface{}{
		"user_id":    event.UserID,
		"event_time": event.EventTime,
		"end_time":   event.EventTime.Add(event.Duration),
	}

	rows, err := r.Postgres.Pool.Query(context.Background(), query, args)
	if err != nil {
		return true, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (r *EventRepo) result(ctx context.Context, id string) (*entity.Event, error) {
	rows, err := r.Postgres.Pool.Query(ctx, "select id, title, descr, user_id, event_time, duration, notification from events where id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := entity.Event{}

	for rows.Next() {
		er := rows.Scan(&res.ID, &res.Title, &res.Desc, &res.UserID, &res.EventTime, &res.Duration, &res.Notification)
		if er != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &res, nil
}
