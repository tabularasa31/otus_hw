package usecase

import (
	"context"
	"fmt"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
)

// EventUseCase -.
type EventUseCase struct {
	repo EventRepo
}

// New -.
func New(r EventRepo) *EventUseCase {
	return &EventUseCase{
		repo: r,
	}
}

func (u *EventUseCase) Create(ctx context.Context, event entity.Event) (*entity.Event, error) {

	eventDB, err := event.Dao()

	res, err := u.repo.CreateEvent(ctx, eventDB)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - CreateEvent - u.repo.CreateEvent: %w", err)
	}

	return res, nil
}

func (u *EventUseCase) Update(ctx context.Context, event entity.Event) (*entity.Event, error) {
	eventDB, err := event.Dao()

	res, err := u.repo.UpdateEvent(ctx, eventDB)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - UpdateEvent - u.repo.UpdateEvent: %w", err)
	}
	return res, nil
}

func (u *EventUseCase) Delete(ctx context.Context, userID int32) error {
	err := u.repo.DeleteEvent(ctx, userID)
	if err != nil {
		return fmt.Errorf("EventUseCase - DeleteEvent - u.repo.DeleteEvent: %w", err)
	}
	return nil
}

func (u *EventUseCase) DailyEvents(ctx context.Context, event entity.Event) ([]entity.Event, error) {
	eventDB, err := event.Dao()

	events, err := u.repo.GetDailyEvents(ctx, eventDB)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - GetDailyEvents - u.repo.GetDailyEvents: %w", err)
	}
	return events, nil
}

func (u *EventUseCase) WeeklyEvents(ctx context.Context, event entity.Event) ([]entity.Event, error) {
	eventDB, err := event.Dao()

	events, err := u.repo.GetWeeklyEvents(ctx, eventDB)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - GetWeeklyEvents - u.repo.GetWeeklyEvents: %w", err)
	}
	return events, nil
}

func (u *EventUseCase) MonthlyEvents(ctx context.Context, event entity.Event) ([]entity.Event, error) {
	eventDB, err := event.Dao()

	events, err := u.repo.GetMonthlyEvents(ctx, eventDB)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - MonthlyEvents - u.repo.MonthlyEvents: %w", err)
	}
	return events, nil
}
