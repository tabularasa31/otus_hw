package entity

import (
	dateconv "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils"
	"time"
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
	date := dateconv.TimeToString(e.StartTime)

	d := dateconv.TimeToString(e.EndTime)

	n := dateconv.TimeToString(e.Notification)

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
