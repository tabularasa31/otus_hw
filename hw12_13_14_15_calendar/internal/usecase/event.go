package usecase

import (
	"context"
	"fmt"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/dateconv"
	"time"

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

func (u *EventUseCase) Create(ctx context.Context, event *entity.Event) (*entity.Event, error) {
	date, err := dateconv.StringToTime(event.StartTime)
	if err != nil {
		return &entity.Event{}, fmt.Errorf("EventUseCase - Create - StringToTime(StartTime): %w", err)
	}

	d, err := dateconv.StringToTime(event.EndTime)
	if err != nil {
		return &entity.Event{}, fmt.Errorf("EventUseCase - Create - StringToTime(EndTime): %w", err)
	}

	n, err := dateconv.StringToTime(event.Notification)
	if err != nil {
		return &entity.Event{}, fmt.Errorf("EventUseCase - Create - StringToTime(Notification): %w", err)
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
		return &entity.Event{}, err
	}

	return result, nil
}

func (u *EventUseCase) Update(ctx context.Context, event entity.Event) (*entity.Event, error) {
	date, err := dateconv.StringToTime(event.StartTime)
	if err != nil {
		return &entity.Event{}, fmt.Errorf("EventUseCase - Update - StringToTime(StartTime): %w", err)
	}

	d, err := dateconv.StringToTime(event.EndTime)
	if err != nil {
		return &entity.Event{}, fmt.Errorf("EventUseCase - Update - StringToTime(EndTime): %w", err)
	}

	n, err := dateconv.StringToTime(event.Notification)
	if err != nil {
		return &entity.Event{}, fmt.Errorf("EventUseCase - Update - StringToTime(Notification): %w", err)
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
		return &entity.Event{}, err
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

func (u *EventUseCase) EventsByDates(ctx context.Context, userID int, start time.Time, end time.Time) ([]entity.Event, error) {
	events, err := u.repo.GetEventsByDates(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - EventsByDates - u.repo.GetEventsByDates: %w", err)
	}
	return events, nil
}

func (u *EventUseCase) EventsByTime(ctx context.Context, start time.Time) ([]entity.Event, error) {
	events, err := u.repo.GetAllEventsByTime(ctx, start)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - EventsByTime - u.repo.GetAllEventsByTime: %w", err)
	}
	return events, nil
}

func (u *EventUseCase) DeleteOldEvents(ctx context.Context, date time.Time) error {
	err := u.repo.DeleteOldEventsFromRepo(ctx, date)
	if err != nil {
		return fmt.Errorf("EventUseCase - DeleteOldEvents - u.repo.DeleteOldEventsFromRepo: %w", err)
	}
	return nil
}
