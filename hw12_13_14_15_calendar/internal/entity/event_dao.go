package entity

import (
	"time"

	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/date_utils"
)

type (
	EventDB struct {
		ID           int32
		Title        string
		Desc         string
		UserID       int
		EventTime    time.Time
		Duration     time.Duration
		Notification time.Duration
	}
)

func (e *EventDB) Dto() *Event {
	date := date_utils.TimeToString(e.EventTime)

	d := date_utils.DurationToString(e.Duration)

	n := date_utils.DurationToString(e.Notification)

	event := Event{
		Title:        e.Title,
		Desc:         e.Desc,
		UserID:       e.UserID,
		EventTime:    date,
		Duration:     d,
		Notification: n,
	}

	return &event
}
