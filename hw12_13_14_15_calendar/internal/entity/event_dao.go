package entity

import (
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/date_utils"
	"time"
)

type (
	EventDB struct {
		Id           int32
		Title        string
		Desc         string
		UserId       int
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
		UserId:       e.UserId,
		EventTime:    date,
		Duration:     d,
		Notification: n,
	}

	return &event
}
