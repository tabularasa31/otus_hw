package memoryrepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	type Case struct {
		name  string
		event entity.Event
		err   error
	}
	t.Run("success event create", func(t *testing.T) {
		userId := int(uuid.New().ID())
		cases := []Case{
			{
				name: "success event create one",
				event: entity.Event{
					Title:        "event 1",
					Desc:         "This is event one",
					UserId:       userId,
					EventTime:    time.Now().Add(time.Hour),
					Duration:     time.Hour,
					Notification: time.Hour,
				},
				err: nil,
			},
			{
				name: "success event create two",
				event: entity.Event{
					Title:        "event 2",
					Desc:         "This is event two",
					UserId:       userId,
					EventTime:    time.Now().Add(3 * time.Hour),
					Duration:     time.Hour,
					Notification: time.Hour * 2,
				},
				err: nil,
			},
		}

		repo := New()
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				err := repo.CreateEvent(context.Background(), c.event)
				require.ErrorIs(t, err, c.err)
			})
		}
	})

	t.Run("invalid event data", func(t *testing.T) {
		userId := int(uuid.New().ID())
		cases := []Case{
			{
				name: "invalid title",
				event: entity.Event{
					Title:        "",
					Desc:         "This is event with empty title",
					UserId:       userId,
					EventTime:    time.Now().Add(3 * time.Hour),
					Duration:     time.Hour,
					Notification: time.Hour * 2,
				},
				err: repo.ErrEventTitle,
			},
			{
				name: "empty time of event",
				event: entity.Event{
					Title:        "Title 333",
					Desc:         "This is event with empty event time",
					UserId:       userId,
					Duration:     time.Hour,
					Notification: time.Hour * 2,
				},
				err: repo.ErrEventTime,
			},
			{
				name: "empty duration",
				event: entity.Event{
					Title:        "Title 222",
					Desc:         "This is event with empty duration",
					UserId:       userId,
					EventTime:    time.Now().Add(5 * time.Hour),
					Notification: time.Hour * 4,
				},
				err: repo.ErrEventDuration,
			},
		}
		stor := New()
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				err := stor.CreateEvent(context.Background(), c.event)
				require.ErrorIs(t, err, c.err)
			})
		}
	})

	t.Run("event time busy", func(t *testing.T) {
		userId := int(uuid.New().ID())
		stor := New()
		err := stor.CreateEvent(context.Background(), entity.Event{
			Title:        "event one",
			Desc:         "event one",
			UserId:       userId,
			EventTime:    time.Date(2022, 12, 30, 15, 0, 0, 0, time.Local),
			Duration:     time.Hour,
			Notification: time.Hour * 2,
		})
		require.NoError(t, err)

		err = stor.CreateEvent(context.Background(), entity.Event{
			Title:        "event two",
			Desc:         "event two",
			UserId:       userId,
			EventTime:    time.Date(2022, 12, 30, 15, 30, 0, 0, time.Local),
			Duration:     time.Hour,
			Notification: time.Hour * 2,
		})
		require.ErrorIs(t, err, repo.ErrEventTimeBusy)

		err = stor.CreateEvent(context.Background(), entity.Event{
			Title:        "event three",
			Desc:         "event three",
			UserId:       userId,
			EventTime:    time.Date(2022, 12, 30, 14, 30, 0, 0, time.Local),
			Duration:     time.Hour,
			Notification: time.Hour * 2,
		})
		require.ErrorIs(t, err, repo.ErrEventTimeBusy)
	})
}
