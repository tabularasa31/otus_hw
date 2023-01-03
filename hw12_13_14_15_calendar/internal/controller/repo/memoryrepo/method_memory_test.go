package memoryrepo_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo/memoryrepo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	type Case struct {
		name  string
		event entity.EventDB
		err   error
	}
	t.Run("success event create", func(t *testing.T) {
		userID := int(uuid.New().ID())
		cases := []Case{
			{
				name: "success event create one",
				event: entity.EventDB{
					Title:        "event 1",
					Desc:         "This is event one",
					UserID:       userID,
					EventTime:    time.Now().Add(time.Hour),
					Duration:     time.Hour,
					Notification: time.Hour,
				},
				err: nil,
			},
			{
				name: "success event create two",
				event: entity.EventDB{
					Title:        "event 2",
					Desc:         "This is event two",
					UserID:       userID,
					EventTime:    time.Now().Add(3 * time.Hour),
					Duration:     time.Hour,
					Notification: time.Hour * 2,
				},
				err: nil,
			},
		}

		eventRepo := memoryrepo.New()
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				_, err := eventRepo.CreateEvent(context.Background(), &c.event)
				require.ErrorIs(t, err, c.err)
			})
		}
	})

	t.Run("event time busy", func(t *testing.T) {
		userID := int(uuid.New().ID())
		eventRepo := memoryrepo.New()
		_, err := eventRepo.CreateEvent(context.Background(), &entity.EventDB{
			Title:        "event one",
			Desc:         "event one",
			UserID:       userID,
			EventTime:    time.Date(2022, 12, 30, 15, 0, 0, 0, time.Local),
			Duration:     time.Hour,
			Notification: time.Hour * 2,
		})
		require.NoError(t, err)

		_, err = eventRepo.CreateEvent(context.Background(), &entity.EventDB{
			Title:        "event two",
			Desc:         "event two",
			UserID:       userID,
			EventTime:    time.Date(2022, 12, 30, 15, 30, 0, 0, time.Local),
			Duration:     time.Hour,
			Notification: time.Hour * 2,
		})
		require.ErrorIs(t, err, repo.ErrEventTimeBusy)

		_, err = eventRepo.CreateEvent(context.Background(), &entity.EventDB{
			Title:        "event three",
			Desc:         "event three",
			UserID:       userID,
			EventTime:    time.Date(2022, 12, 30, 14, 30, 0, 0, time.Local),
			Duration:     time.Hour,
			Notification: time.Hour * 2,
		})
		require.ErrorIs(t, err, repo.ErrEventTimeBusy)
	})
}
