package entity

import (
	"time"
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

//ID - уникальный идентификатор события (можно воспользоваться UUID);
//Заголовок - короткий текст;
//Дата и время события;
//Длительность события (или дата и время окончания);
//Описание события - длинный текст, опционально;
//ID пользователя, владельца события;
//За сколько времени высылать уведомление, опционально.
