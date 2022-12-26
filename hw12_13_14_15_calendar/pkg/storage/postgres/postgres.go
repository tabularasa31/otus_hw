package postgres

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

//
//type Postgres struct {
//	db  *sql.DB
//	dsn string
//}

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool *pgxpool.Pool
}

func New(url string) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

// Close -.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

//
//func (s *Postgres) CreateEvent(ctx context.Context, event models.Event) error {
//	if err := s.EventValidate(event); err != nil {
//		return err
//	}
//
//	check, er := s.IsEventTimeBusy(event)
//	if !check {
//		return errapp.ErrEventTimeBusy
//	}
//	if er != nil {
//		return er
//	}
//
//	query := `INSERT INTO events(userid, title, descr, event_time, duration)
//			  VALUES (:userId, :title, :descr, :event_time, :duration)`
//
//	args := map[string]interface{}{
//		"userId":     event.UserId,
//		"title":      event.Title,
//		"descr":      event.Desc,
//		"event_time": event.EventTime,
//		"duration":   event.Duration,
//	}
//	_, e := s.db.QueryContext(ctx, query, args)
//	if e != nil {
//		return fmt.Errorf("failed to create: %w", e)
//	}
//
//	return nil
//}

//func (s *Postgres) UpdateEvent(ctx context.Context, event models.Event) error {
//	if err := s.EventValidate(event); err != nil {
//		return err
//	}
//
//	check, er := s.IsEventTimeBusy(event)
//	if !check {
//		return errapp.ErrEventTimeBusy
//	}
//	if er != nil {
//		return er
//	}
//
//	query := `UPDATE events
//				SET userid = :userId,
//				    title = :title,
//				    descr = :descr,
//				    event_time = :event_time,
//				    duration = :duration,
//				    notification = :notification
//				WHERE id = :id`
//
//	args := map[string]interface{}{
//		"id":           event.Id,
//		"userId":       event.UserId,
//		"title":        event.Title,
//		"descr":        event.Desc,
//		"event_time":   event.EventTime,
//		"duration":     event.Duration,
//		"notification": event.Notification,
//	}
//	_, e := s.db.ExecContext(ctx, query, args)
//	if e != nil {
//		return fmt.Errorf("failed to update: %w", e)
//	}
//
//	return nil
//}
//
//// DeleteEvent Удалить (ID события);
//func (s *Postgres) DeleteEvent(ctx context.Context, Id int32) error {
//	_, err := s.db.Exec(`delete from events where id = $1`, Id)
//	return err
//}
//
//// GetDailyEvents СписокСобытийНаДень (дата);
//// Выводит все события, которые начинаются в заданный день
//func (s *Postgres) GetDailyEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
//	var events []models.Event
//	query := `SELECT id, title, descr, user_id, event_time, duration, notification
//       			FROM events
//       			WHERE DATE_TRUNC('day', event_time) = :date`
//	args := map[string]interface{}{
//		"date": date.Day(),
//	}
//	rows, err := s.db.QueryContext(ctx, query, args)
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		if closeErr := rows.Close(); closeErr != nil {
//			fmt.Errorf("failed to close row: %w", closeErr)
//		}
//	}()
//
//	for rows.Next() {
//		var event models.Event
//		if err := rows.Scan(&event.Id, &event.Title, &event.Desc,
//			&event.UserId, &event.EventTime, &event.Duration, &event.Notification); err != nil {
//			return events, err
//		}
//		events = append(events, event)
//	}
//
//	return events, nil
//}
//
//// GetWeeklyEvents СписокСобытийНаНеделю (дата начала недели);
//// Выводит список событий за 7 дней, начиная с дня начала
//func (s *Postgres) GetWeeklyEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
//	var events []models.Event
//	query := `SELECT id, title, descr, user_id, event_time, duration, notification
//       			FROM events
//       			WHERE DATE_TRUNC('week', event_time)
//       			BETWEEN date :date AND date :date + 7 * interval '1 day'`
//	args := map[string]interface{}{
//		"date": date.Day(),
//	}
//	rows, err := s.db.QueryContext(ctx, query, args)
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		if closeErr := rows.Close(); closeErr != nil {
//			fmt.Errorf("failed to close row: %w", closeErr)
//		}
//	}()
//	for rows.Next() {
//		var event models.Event
//		if err := rows.Scan(&event.Id, &event.Title, &event.Desc,
//			&event.UserId, &event.EventTime, &event.Duration, &event.Notification); err != nil {
//			return events, err
//		}
//		events = append(events, event)
//	}
//	return events, nil
//}
//
//// GetMonthlyEvents СписокСобытийНaМесяц (дата начала месяца)
//// Выводит список событий за 30 дней, начиная с дня начала
//func (s *Postgres) GetMonthlyEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
//	var events []models.Event
//	query := `SELECT id, title, descr, user_id, event_time, duration, notification
//       			FROM events
//       			WHERE DATE_TRUNC('week', event_time)
//       			BETWEEN date :date AND date :date + 30 * interval '1 day'`
//	args := map[string]interface{}{
//		"date": date.Day(),
//	}
//	rows, err := s.db.QueryContext(ctx, query, args)
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		if closeErr := rows.Close(); closeErr != nil {
//			fmt.Errorf("failed to close row: %w", closeErr)
//		}
//	}()
//
//	for rows.Next() {
//		var event models.Event
//		if err := rows.Scan(&event.Id, &event.Title, &event.Desc,
//			&event.UserId, &event.EventTime, &event.Duration, &event.Notification); err != nil {
//			return events, err
//		}
//		events = append(events, event)
//	}
//	return events, nil
//}
//
//// IsEventTimeBusy проверка на занятость времени
//func (s *Postgres) IsEventTimeBusy(event models.Event) (bool, error) {
//	//TODO: Написать проверку времени на занятость
//	query := `SELECT id
//			  FROM events
//			  WHERE userid = :userId
//			    AND event_time BETWEEN :event_time AND :end_time
//			  LIMIT 1`
//	args := map[string]interface{}{
//		"userId":     event.UserId,
//		"event_time": event.EventTime,
//		"end_time":   event.EventTime.Add(event.Duration),
//	}
//
//	rows, err := s.db.QueryContext(context.Background(), query, args)
//	if err != nil {
//		return true, err
//	}
//	defer func() {
//		if closeErr := rows.Close(); closeErr != nil {
//			fmt.Errorf("failed to close row: %w", closeErr)
//		}
//	}()
//
//	return rows.Next(), nil
//}
//
//// EventValidate проверка ивента на валидность
//func (s *Postgres) EventValidate(event models.Event) error {
//	//TODO написать ивент валидатор
//	_ = event
//	return nil
//}
