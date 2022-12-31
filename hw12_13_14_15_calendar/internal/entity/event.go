package entity

import (
	"time"
)

type Event struct {
	Id           int32         `json:"id"`
	Title        string        `json:"title"`
	Desc         string        `json:"desc"`
	UserId       int           `json:"user_id"`
	EventTime    time.Time     `json:"event_time"`
	Duration     time.Duration `json:"duration"`
	Notification time.Duration `json:"notification"`
}

//ID - уникальный идентификатор события (можно воспользоваться UUID);
//Заголовок - короткий текст;
//Дата и время события;
//Длительность события (или дата и время окончания);
//Описание события - длинный текст, опционально;
//ID пользователя, владельца события;
//За сколько времени высылать уведомление, опционально.
