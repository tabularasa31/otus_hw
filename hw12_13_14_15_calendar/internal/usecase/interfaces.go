package usecase

import (
	"context"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"time"
)

type (
	Event interface {
		Create(context.Context, entity.Event) error
		Update(context.Context, entity.Event) error
		Delete(context.Context, int32) error
		DailyEvents(context.Context, time.Time) ([]entity.Event, error)
		WeeklyEvents(context.Context, time.Time) ([]entity.Event, error)
		MonthlyEvents(context.Context, time.Time) ([]entity.Event, error)
	}

	EventRepo interface {
		CreateEvent(context.Context, entity.Event) error
		UpdateEvent(context.Context, entity.Event) error
		DeleteEvent(context.Context, int32) error
		GetDailyEvents(context.Context, time.Time) ([]entity.Event, error)
		GetWeeklyEvents(context.Context, time.Time) ([]entity.Event, error)
		GetMonthlyEvents(context.Context, time.Time) ([]entity.Event, error)
	}
)
