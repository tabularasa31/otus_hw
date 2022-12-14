package interfaces

import (
	"context"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/domain/models"
	"time"
)

type Storage interface {
	Create(ctx context.Context, event models.Event) error
	Update(ctx context.Context, event models.Event) error
	Delete(Id int32) error
	GetDailyEvents(ctx context.Context, date time.Time) ([]models.Event, error)
	GetWeeklyEvents(ctx context.Context, date time.Time) ([]models.Event, error)
	GetMonthlyEvents(ctx context.Context, date time.Time) ([]models.Event, error)
}
