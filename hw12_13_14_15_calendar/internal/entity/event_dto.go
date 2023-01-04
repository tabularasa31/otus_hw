package entity

import (
	"fmt"
	errapp "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/date_utils"
)

type (
	Event struct {
		ID           int32  `json:"id"`
		Title        string `json:"title"`
		Desc         string `json:"desc"`
		UserID       int    `json:"user_id"`
		EventTime    string `json:"event_time"`
		Duration     string `json:"duration"`
		Notification string `json:"notification"`
	}
)

func (e *Event) Dao() (*EventDB, error) {
	//if err := e.eventValidate(); err != nil {
	//	return nil, fmt.Errorf("error event validation: %w", err)
	//}

	date, err := date_utils.StringToTime(e.EventTime)
	if err != nil {
		return nil, fmt.Errorf("StringToTime: %w", err)
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
		UserID:       e.UserID,
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
	case e.UserID == 0:
		return errapp.ErrEventUserID
	case e.EventTime == "":
		return errapp.ErrEventTime
	case e.Duration == "":
		return errapp.ErrEventDuration
	}
	return nil
}
