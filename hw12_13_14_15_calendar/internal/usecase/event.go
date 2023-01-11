package usecase

import (
	"context"
	"fmt"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/date_utils"
	"time"
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
	date, err := date_utils.StringToTime(event.StartTime)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - Create - StringToTime(StartTime): %w", err)
	}

	d, err := date_utils.StringToTime(event.EndTime)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - Create - StringToTime(EndTime): %w", err)
	}

	n, err := date_utils.StringToTime(event.Notification)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - Create - StringToTime(Notification): %w", err)
	}

	eventDB := entity.EventDB{
		Title:        event.Title,
		Desc:         event.Desc,
		UserID:       event.UserID,
		StartTime:    date,
		EndTime:      d,
		Notification: n,
	}

	result, err := u.repo.CreateEvent(ctx, &eventDB)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *EventUseCase) Update(ctx context.Context, event entity.Event) (*entity.Event, error) {
	date, err := date_utils.StringToTime(event.StartTime)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - Update - StringToTime(StartTime): %w", err)
	}

	d, err := date_utils.StringToTime(event.EndTime)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - Update - StringToTime(EndTime): %w", err)
	}

	n, err := date_utils.StringToTime(event.Notification)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - Update - StringToTime(Notification): %w", err)
	}

	eventDB := entity.EventDB{
		ID:           event.ID,
		Title:        event.Title,
		Desc:         event.Desc,
		UserID:       event.UserID,
		StartTime:    date,
		EndTime:      d,
		Notification: n,
	}

	res, err := u.repo.UpdateEvent(ctx, &eventDB)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *EventUseCase) Delete(ctx context.Context, userID int) error {
	err := u.repo.DeleteEvent(ctx, userID)
	if err != nil {
		return fmt.Errorf("EventUseCase - DeleteEvent - u.repo.DeleteEvent: %w", err)
	}
	return nil
}

func (u *EventUseCase) DailyEvents(ctx context.Context, userID int, date time.Time) ([]entity.Event, error) {
	events, err := u.repo.GetDailyEvents(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - DailyEvents - u.repo.GetDailyEvents: %w", err)
	}
	return events, nil
}

func (u *EventUseCase) WeeklyEvents(ctx context.Context, userID int, date time.Time) ([]entity.Event, error) {
	events, err := u.repo.GetWeeklyEvents(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - WeeklyEvents - u.repo.GetWeeklyEvents: %w", err)
	}
	return events, nil
}

func (u *EventUseCase) MonthlyEvents(ctx context.Context, userID int, date time.Time) ([]entity.Event, error) {
	events, err := u.repo.GetMonthlyEvents(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - MonthlyEvents - u.repo.MonthlyEvents: %w", err)
	}
	return events, nil
}
