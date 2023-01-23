package usecase

import (
	"context"
	"time"

	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
)

type (
	EventRepo interface {
		CreateEvent(context.Context, *entity.EventDB) (*entity.Event, error)
		UpdateEvent(context.Context, *entity.EventDB) (*entity.Event, error)
		DeleteEvent(context.Context, int) error
		GetEventsByDates(context.Context, int, time.Time, time.Time) ([]entity.Event, error)
	}
)
