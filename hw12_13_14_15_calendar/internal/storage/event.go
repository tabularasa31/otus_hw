package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidEventTitle  = errors.New("invalid event title")
	ErrInvalidEventTime   = errors.New("invalid event time")
	ErrInvalidEventUserID = errors.New("invalid event user id")
)

type Event struct {
	Id           int32
	Title        string
	Desc         string
	UserId       int
	EventTime    time.Time
	Duration     time.Duration
	Notification time.Time
}

type Storage interface {
	Create(ctx context.Context, event Event) error
	Update(ctx context.Context, event Event) error
	Delete(Id int32) error
	GetDailyEvents(ctx context.Context, date time.Time) ([]Event, error)
	GetWeeklyEvents(ctx context.Context, date time.Time) ([]Event, error)
	GetMonthlyEvents(ctx context.Context, date time.Time) ([]Event, error)
}

//ID - уникальный идентификатор события (можно воспользоваться UUID);
//Заголовок - короткий текст;
//Дата и время события;
//Длительность события (или дата и время окончания);
//Описание события - длинный текст, опционально;
//ID пользователя, владельца события;
//За сколько времени высылать уведомление, опционально.
