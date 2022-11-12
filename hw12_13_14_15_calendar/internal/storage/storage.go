package storage

import (
	"context"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/storage/sql"
	"time"
)

type Storage interface {
	Create(ctx context.Context, event Event) error
	Update(ctx context.Context, event Event) error
	Delete(Id int32) error
	GetDailyEvents(ctx context.Context, date time.Time) ([]Event, error)
	GetWeeklyEvents(ctx context.Context, date time.Time) ([]Event, error)
	GetMonthlyEvents(ctx context.Context, date time.Time) ([]Event, error)
}

func New(StorageConf config.StorageConf) Storage {
	switch StorageConf.Type {
	case "memory":
		return memorystorage.New()
	case "sql":
		return sqlstorage.New(StorageConf.Dsn)
	default:
		return memorystorage.New()
	}
}
