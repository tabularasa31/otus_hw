package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/utils"
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
	date, err := utils.StringToTime(event.StartTime)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - Create - StringToTime(StartTime): %w", err)
	}

	d, err := utils.StringToTime(event.EndTime)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - Create - StringToTime(EndTime): %w", err)
	}

	n, err := utils.StringToTime(event.Notification)
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
	date, err := utils.StringToTime(event.StartTime)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - Update - StringToTime(StartTime): %w", err)
	}

	d, err := utils.StringToTime(event.EndTime)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - Update - StringToTime(EndTime): %w", err)
	}

	n, err := utils.StringToTime(event.Notification)
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
	fmt.Println("------Starting DeleteOldEventsFromRepo --------")
	err := u.repo.DeleteOldEventsFromRepo(ctx, date)
	if err != nil {
		return fmt.Errorf("EventUseCase - DeleteOldEvents - u.repo.DeleteOldEventsFromRepo: %w", err)
	}
	return nil
}
