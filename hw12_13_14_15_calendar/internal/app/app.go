package app

import (
	"context"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/domain/interfaces"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/domain/models"
)

type App struct {
	storage interfaces.Storage
	logger  Logger
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

func New(logger Logger, storage interfaces.Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) Create(ctx context.Context, event models.Event) error {
	return a.storage.Create(ctx, event)
}

func (a *App) Update(ctx context.Context, event models.Event) error {
	return a.storage.Update(ctx, event)
}

// TODO
