package usecase

import (
	"context"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"time"
)

type Storage interface {
	CreateEvent(ctx context.Context, event entity.Event) error
	UpdateEvent(ctx context.Context, event entity.Event) error
	DeleteEvent(ctx context.Context, Id int32) error
	GetDailyEvents(ctx context.Context, date time.Time) ([]entity.Event, error)
	GetWeeklyEvents(ctx context.Context, date time.Time) ([]entity.Event, error)
	GetMonthlyEvents(ctx context.Context, date time.Time) ([]entity.Event, error)
}
