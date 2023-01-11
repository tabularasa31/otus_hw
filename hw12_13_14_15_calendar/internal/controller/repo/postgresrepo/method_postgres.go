package postgresrepo

import (
	"context"
	"fmt"
	errapp "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/storage/postgres"
	"time"
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
	if check {
		return nil, errapp.ErrEventTimeBusy
	}
	if er != nil {
		return nil, fmt.Errorf("postgres - CreateEvent - r.isEventTimeBusy: %w", er)
	}

	sql, args, err := r.Builder.
		Insert("events").
		Columns("user_id, title, descr, start_time, end_time, notification").
		Values(eventDB.UserID, eventDB.Title, eventDB.Desc, eventDB.StartTime, eventDB.EndTime, eventDB.Notification).
		Suffix("returning id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("postgres - CreateEvent - r.Builder: %w", err)
	}

	var lastInsertID int
	err = r.Postgres.Pool.QueryRow(ctx, sql, args...).Scan(&lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("postgres - CreateEvent - r.Postgres.Pool.QueryRow: %w", err)
	}

	result, err := r.result(ctx, lastInsertID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateEvent Обновить (событие).
func (r *EventRepo) UpdateEvent(ctx context.Context, eventDB *entity.EventDB) (*entity.Event, error) {
	check, er := r.isEventTimeBusy(*eventDB)
	if check {
		return nil, errapp.ErrEventTimeBusy
	}
	if er != nil {
		return nil, fmt.Errorf("postgres - UpdateEvent - r.isEventTimeBusy: %w", er)
	}

	sql, args, err := r.Builder.
		Update("events").
		Set("user_id", eventDB.UserID).
		Set("title", eventDB.Title).
		Set("descr", eventDB.Desc).
		Set("start_time", eventDB.StartTime).
		Set("end_time", eventDB.EndTime).
		Set("notification", eventDB.Notification).
		Where("id=?", eventDB.ID).ToSql()
	if err != nil {
		return nil, fmt.Errorf("postgres - UpdateEvent - r.Builder: %w", err)
	}

	fmt.Println("QUERY ------- ", sql)
	fmt.Println("ARGS ------- ", args)

	qqq, e := r.Postgres.Pool.Exec(ctx, sql, args...)
	if e != nil {
		return nil, fmt.Errorf("postgres - UpdateEvent - r.Postgres.Pool.Exec: %w", e)
	}

	fmt.Println("-----------QQQ ------- ", qqq.RowsAffected(), qqq.Update())

	fmt.Println("-----------EVENT ID------------, ", eventDB.ID)

	result, err := r.result(ctx, eventDB.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteEvent Удалить (ID события).
func (r *EventRepo) DeleteEvent(ctx context.Context, id int) error {
	_, err := r.Postgres.Pool.Exec(ctx, `delete from events where id = $1`, id)
	return err
}

// GetDailyEvents СписокСобытийНаДень (дата).
// Выводит все события, которые начинаются в заданный день.
func (r *EventRepo) GetDailyEvents(ctx context.Context, userID int, date time.Time) ([]entity.Event, error) {
	var events []entity.Event

	query := `SELECT id, title, descr, user_id, start_time, end_time, notification
       			FROM events 
       			WHERE user_id = $1 
       			  AND DATE_TRUNC('day', start_time) = $2`

	rows, err := r.Postgres.Pool.Query(ctx, query, userID, date.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("GetDailyEvents - r.Postgres.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var eventDB entity.EventDB
		if er := rows.Scan(&eventDB.ID, &eventDB.Title, &eventDB.Desc, &eventDB.UserID, &eventDB.StartTime, &eventDB.EndTime, &eventDB.Notification); er != nil {
			return events, fmt.Errorf("GetDailyEvents - rows.Scan: %w", er)
		}
		events = append(events, *eventDB.Dto())
	}
	return events, nil
}

// GetWeeklyEvents СписокСобытийНаНеделю (дата начала недели).
// Выводит список событий за 7 дней, начиная с дня начала.
func (r *EventRepo) GetWeeklyEvents(ctx context.Context, userID int, date time.Time) ([]entity.Event, error) {
	var events []entity.Event

	startDate := date
	endDate := startDate.Add(7 * 24 * time.Hour)

	events, err := r.getEventsByDates(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// GetMonthlyEvents СписокСобытийНaМесяц (дата начала месяца).
// Выводит список событий за 30 дней, начиная с дня начала.
func (r *EventRepo) GetMonthlyEvents(ctx context.Context, userID int, date time.Time) ([]entity.Event, error) {
	var events []entity.Event

	startDate := date
	endDate := startDate.Add(30 * 24 * time.Hour)

	events, err := r.getEventsByDates(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// isEventTimeBusy проверка на занятость времени.
func (r *EventRepo) isEventTimeBusy(eventDB entity.EventDB) (bool, error) {
	query := `SELECT id 
			  FROM events 
			  WHERE user_id = $1 
			    AND start_time BETWEEN $2 AND $3
			  LIMIT 1`

	rows, err := r.Postgres.Pool.Query(context.Background(), query, eventDB.UserID, eventDB.StartTime, eventDB.EndTime)
	if err != nil {
		return true, fmt.Errorf("postgres - isEventTimeBusy - r.Postgres.Pool.Query: %w", err)
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (r *EventRepo) result(ctx context.Context, id int) (*entity.Event, error) {

	fmt.Println("--------ID----------  ", id)

	rows, err := r.Postgres.Pool.Query(ctx, "select id, title, descr, user_id, start_time, end_time, notification from events where id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("postgres - result - r.Postgres.Pool.Query: %w", err)
	}
	defer rows.Close()

	res := entity.EventDB{}

	for rows.Next() {
		er := rows.Scan(&res.ID, &res.Title, &res.Desc, &res.UserID, &res.StartTime, &res.EndTime, &res.Notification)
		if er != nil {
			return nil, fmt.Errorf("postgres - result - rows.Scan: %w", err)
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("postgres - result - rows.Err: %w", err)
	}

	fmt.Println("RESULT RESULT -----------", res)

	return res.Dto(), nil
}

func (r *EventRepo) getEventsByDates(ctx context.Context, uid int, startDate time.Time, endDate time.Time) ([]entity.Event, error) {
	var events []entity.Event
	query := `SELECT id, title, descr, start_time, end_time
       			FROM events 
       			WHERE user_id = :uid
       			AND start_time
       			BETWEEN date :startDate AND date :endDate`
	args := map[string]interface{}{
		"uid":       uid,
		"startDate": startDate,
		"endDate":   endDate,
	}
	rows, err := r.Postgres.Pool.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("postgres - getEventsByDates - r.Postgres.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var event entity.EventDB
		if e := rows.Scan(&event.ID, &event.Title, &event.Desc, &event.StartTime, &event.EndTime); e != nil {
			return events, fmt.Errorf("postgres - getEventsByDates - rows.Scan: %w", e)
		}
		events = append(events, *event.Dto())
	}

	return events, nil
}
