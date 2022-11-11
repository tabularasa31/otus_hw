package storage

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidEventTitle  = errors.New("invalid event title")
	ErrInvalidEventTime   = errors.New("invalid event time")
	ErrInvalidEventUserID = errors.New("invalid event user id")
)

type Event struct {
	Id           uuid.UUID
	Title        string
	EventTime    time.Time
	Duration     time.Duration
	Desc         string
	UserId       int
	Notification time.Duration
}

//ID - уникальный идентификатор события (можно воспользоваться UUID);
//Заголовок - короткий текст;
//Дата и время события;
//Длительность события (или дата и время окончания);
//Описание события - длинный текст, опционально;
//ID пользователя, владельца события;
//За сколько времени высылать уведомление, опционально.

func (event Event) EventValidate() error {
	//TODO написать ивент валидатор
	return nil
}
