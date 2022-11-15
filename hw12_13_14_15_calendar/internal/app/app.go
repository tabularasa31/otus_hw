package app

import (
	"context"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type App struct {
	storage Storage
	logger  Logger
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type Storage interface {
	Create(ctx context.Context, event storage.Event) error
	Update(ctx context.Context, event storage.Event) error
	Delete(Id int32) error
	GetDailyEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetWeeklyEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetMonthlyEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) Create(ctx context.Context, event storage.Event) error {
	return a.storage.Create(ctx, event)
}

func (a *App) Update(ctx context.Context, event storage.Event) error {
	return a.storage.Update(ctx, event)
}

// TODO
