package entity

import (
	"fmt"
	errapp "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/date_utils"
)

type (
	Event struct {
		ID           int32  `json:"id" example:"auto"`
		Title        string `json:"title" binding:"required" example:"title"`
		Desc         string `json:"desc" binding:"required" example:"description"`
		UserID       int    `json:"userId" binding:"required" example:"user id"`
		EventTime    string `json:"eventTime" binding:"required" example:"start event time in RFC3339 format"`
		Duration     string `json:"duration" binding:"required" example:"event duration in time.Duration format"`
		Notification string `json:"notification" binding:"required" example:"event notification in time.Duration format"`
	}
)

func (e *Event) Dao() (*EventDB, error) {
	if err := e.eventValidate(); err != nil {
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
