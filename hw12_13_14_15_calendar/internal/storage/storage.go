package storage

import (
	"github.com/google/uuid"
	memorystorage "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/storage/sql"
)

type Storage interface {
	Create(event Event) error
	Update(event Event) error
	Delete(Id uuid.UUID) error
	GetDailyEvents(date string) ([]Event, error)
	GetWeeklyEvents(date string) ([]Event, error)
	GetMonthlyEvents(date string) ([]Event, error)
}

func New(Type string) Storage {
	switch Type {
	case "memory":
		return memorystorage.New()
	case "sql":
		return sqlstorage.New()
	default:
		return memorystorage.New()
	}
}
