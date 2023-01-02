package usecase

import (
	"context"

	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
)

type (
	EventRepo interface {
		CreateEvent(context.Context, *entity.EventDB) (*entity.Event, error)
		UpdateEvent(context.Context, *entity.EventDB) (*entity.Event, error)
		DeleteEvent(context.Context, int32) error
		GetDailyEvents(context.Context, *entity.EventDB) ([]entity.Event, error)
		GetWeeklyEvents(context.Context, *entity.EventDB) ([]entity.Event, error)
		GetMonthlyEvents(context.Context, *entity.EventDB) ([]entity.Event, error)
	}
)
