package memoryrepo_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo/memoryrepo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
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
					StartTime:    time.Now().Add(time.Hour),
					EndTime:      time.Now().Add(2 * time.Hour),
					Notification: time.Now(),
				},
				err: nil,
			},
			{
				name: "success event create two",
				event: entity.EventDB{
					Title:        "event 2",
					Desc:         "This is event two",
					UserID:       userID,
					StartTime:    time.Now().Add(3 * time.Hour),
					EndTime:      time.Now().Add(4 * time.Hour),
					Notification: time.Now(),
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
}
