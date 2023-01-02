package entity

import (
	"fmt"
	errapp "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/usecase"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/date_utils"
)

type (
	Event struct {
		Id           int32  `json:"id" example:"auto"`
		Title        string `json:"title" binding:"required" example:"title"`
		Desc         string `json:"desc" binding:"required" example:"description"`
		UserId       int    `json:"user_id" binding:"required" example:"user id"`
		EventTime    string `json:"event_time" binding:"required" example:"start event time in RFC3339 format"`
		Duration     string `json:"duration" binding:"required" example:"event duration in time.Duration format"`
		Notification string `json:"notification" binding:"required" example:"за сколько времени до начала события прислать уведомление in time.Duration format"`
	}
)

//ID - уникальный идентификатор события (можно воспользоваться UUID);
//Заголовок - короткий текст;
//Дата и время события;
//Длительность события (или дата и время окончания);
//Описание события - длинный текст, опционально;
//ID пользователя, владельца события;
//За сколько времени высылать уведомление, опционально.

func (e *Event) Dao() (*EventDB, error) {
	if err := e.eventValidate; err != nil {
		return nil, fmt.Errorf("error event validation: %w", err)
	}

	date, err := date_utils.StringToTime(e.EventTime)
	if err != nil {
		return nil, err
	}

	d, err := date_utils.StringToDuration(e.Duration)
	if err != nil {
		return nil, err
	}

	n, err := date_utils.StringToDuration(e.Notification)
	if err != nil {
		return nil, err
	}

	eventDB := EventDB{
		Title:        e.Title,
		Desc:         e.Desc,
		UserId:       e.UserId,
		EventTime:    date,
		Duration:     d,
		Notification: n,
	}

	return &eventDB, nil
}

// eventValidate проверка ивента на валидность
func (e *Event) eventValidate() error {
	switch {
	case e.Title == "":
		return errapp.ErrEventTitle
	case e.UserId == 0:
		return errapp.ErrEventUserID
	case e.EventTime == "":
		return errapp.ErrEventTime
	case e.Duration == "":
		return errapp.ErrEventDuration
	}
	return nil
}
