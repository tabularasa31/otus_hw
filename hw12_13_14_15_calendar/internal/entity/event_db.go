package entity

import (
	"time"

	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/utils"
)

type (
	EventDB struct {
		ID           int
		Title        string
		Desc         string
		UserID       int
		StartTime    time.Time
		EndTime      time.Time
		Notification time.Time
	}
)

func (e *EventDB) Dto() *Event {
	date := utils.TimeToString(e.StartTime)

	d := utils.TimeToString(e.EndTime)

	n := utils.TimeToString(e.Notification)

	event := Event{
		ID:           e.ID,
		Title:        e.Title,
		Desc:         e.Desc,
		UserID:       e.UserID,
		StartTime:    date,
		EndTime:      d,
		Notification: n,
	}

	return &event
}
