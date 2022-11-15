package memorystorage_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/storage/memory"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	type Case struct {
		name  string
		event storage.Event
		err   error
	}
	t.Run("success event create", func(t *testing.T) {
		userId := int(uuid.New().ID())
		cases := []Case{
			{
				name: "success event create one",
				event: storage.Event{
					Title:        "event 1",
					Desc:         "This is event one",
					UserId:       userId,
					EventTime:    time.Now().Add(time.Hour),
					Duration:     time.Hour,
					Notification: time.Now(),
				},
				err: nil,
			},
			{
				name: "success event create two",
				event: storage.Event{
					Title:        "event 2",
					Desc:         "This is event two",
					UserId:       userId,
					EventTime:    time.Now().Add(3 * time.Hour),
					Duration:     time.Hour,
					Notification: time.Now().Add(2 * time.Hour),
				},
				err: nil,
			},
		}

		stor := memorystorage.New()
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				err := stor.Create(context.Background(), c.event)
				require.ErrorIs(t, err, c.err)
			})
		}
	})

	t.Run("invalid event data", func(t *testing.T) {
		userId := int(uuid.New().ID())
		cases := []Case{
			{
				name: "invalid title",
				event: storage.Event{
					Title:        "",
					Desc:         "This is event with empty title",
					UserId:       userId,
					EventTime:    time.Now().Add(3 * time.Hour),
					Duration:     time.Hour,
					Notification: time.Now().Add(2 * time.Hour),
				},
				err: storage.ErrInvalidEventTitle,
			},
			{
				name: "empty time of event",
				event: storage.Event{
					Title:        "Title 333",
					Desc:         "This is event with empty event time",
					UserId:       userId,
					Duration:     time.Hour,
					Notification: time.Now().Add(2 * time.Hour),
				},
				err: storage.ErrInvalidEventTime,
			},
			{
				name: "empty duration",
				event: storage.Event{
					Title:        "Title 222",
					Desc:         "This is event with empty duration",
					UserId:       userId,
					EventTime:    time.Now().Add(5 * time.Hour),
					Notification: time.Now().Add(4 * time.Hour),
				},
				err: storage.ErrInvalidEventTime,
			},
		}
		stor := memorystorage.New()
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				err := stor.Create(context.Background(), c.event)
				require.ErrorIs(t, err, c.err)
			})
		}
	})

	t.Run("event time busy", func(t *testing.T) {
		userId := int(uuid.New().ID())
		stor := memorystorage.New()
		err := stor.Create(context.Background(), storage.Event{
			Title:        "event one",
			Desc:         "event one",
			UserId:       userId,
			EventTime:    time.Date(2022, 12, 30, 15, 0, 0, 0, time.Local),
			Duration:     time.Hour,
			Notification: time.Now().Add(2 * time.Hour),
		})
		require.NoError(t, err)

		err = stor.Create(context.Background(), storage.Event{
			Title:        "event two",
			Desc:         "event two",
			UserId:       userId,
			EventTime:    time.Date(2022, 12, 30, 15, 30, 0, 0, time.Local),
			Duration:     time.Hour,
			Notification: time.Now().Add(2 * time.Hour),
		})
		require.ErrorIs(t, err, memorystorage.ErrEventTimeBusy)

		err = stor.Create(context.Background(), storage.Event{
			Title:        "event three",
			Desc:         "event three",
			UserId:       userId,
			EventTime:    time.Date(2022, 12, 30, 14, 30, 0, 0, time.Local),
			Duration:     time.Hour,
			Notification: time.Now().Add(2 * time.Hour),
		})
		require.ErrorIs(t, err, memorystorage.ErrEventTimeBusy)
	})
}
