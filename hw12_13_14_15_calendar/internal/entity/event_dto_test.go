package entity

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo"
	"testing"
)

func TestCreate(t *testing.T) {
	type Case struct {
		name  string
		event Event
		err   error
	}

	t.Run("invalid event data", func(t *testing.T) {
		userId := int(uuid.New().ID())
		cases := []Case{
			{
				name: "invalid title",
				event: Event{
					Desc:         "This is event with empty title",
					UserId:       userId,
					EventTime:    "2006-01-02 15:04:05",
					Duration:     "1h",
					Notification: "2h",
				},
				err: repo.ErrEventTitle,
			},
			{
				name: "empty time of event",
				event: Event{
					Title:        "This is title",
					Desc:         "This is event with empty event time",
					UserId:       userId,
					Duration:     "1h",
					Notification: "2h",
				},
				err: repo.ErrEventTime,
			},
			{
				name: "empty duration",
				event: Event{
					Title:        "This is title",
					Desc:         "This is event with empty duration",
					UserId:       userId,
					EventTime:    "2006-01-02 15:04:05",
					Notification: "30min",
				},
				err: repo.ErrEventDuration,
			},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				err := c.event.eventValidate()
				require.Error(t, err, c.err)
			})
		}
	})
}
