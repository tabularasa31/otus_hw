package memoryrepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/usecase"
	errapp "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/usecase"
	"sync"
	"time"
)

// EventRepo -.
type EventRepo struct {
	events map[int]map[int32]entity.Event
	mu     sync.RWMutex
}

// New -.
func New() *EventRepo {
	m := sync.RWMutex{}
	events := make(map[int]map[int32]entity.Event)
	return &EventRepo{
		events: events,
		mu:     m,
	}
}

// CreateEvent Создать (событие)
func (r *EventRepo) CreateEvent(ctx context.Context, event entity.Event) error {
	if err := r.eventValidate(event); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// CreateEvent unique ID
	event.Id = int32(uuid.New().ID())

	userEvents, ok := r.events[event.UserId]
	if !ok {
		r.events[event.UserId] = make(map[int32]entity.Event)
	}

	if !r.isEventTimeBusy(userEvents, event) {
		return usecase.ErrEventTimeBusy
	}

	r.events[event.UserId][event.Id] = event
	return nil
}

// UpdateEvent Обновить (ID пользователя, ID события, событие);
func (r *EventRepo) UpdateEvent(ctx context.Context, event entity.Event) error {
	if err := r.eventValidate(event); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	userEvents, ok := r.events[event.UserId]
	if !ok {
		return errapp.ErrEventNotFound
	}

	updatedEvent, ok := userEvents[event.Id]
	if !ok {
		return errapp.ErrEventNotFound
	}

	if !r.isEventTimeBusy(userEvents, event) {
		return errapp.ErrEventTimeBusy
	}

	updatedEvent.Title = event.Title
	updatedEvent.Desc = event.Desc
	updatedEvent.EventTime = event.EventTime
	updatedEvent.Duration = event.Duration

	r.events[event.UserId][event.Id] = updatedEvent
	return nil
}

// DeleteEvent Удалить (ID события);
func (r *EventRepo) DeleteEvent(ctx context.Context, Id int32) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, userEvents := range r.events {
		if _, ok := userEvents[Id]; !ok {
			delete(userEvents, Id)
			return nil
		}
	}
	return usecase.ErrEventNotFound
}

// GetDailyEvents СписокСобытийНаДень (дата);
// Выводит все события, которые начинаются в заданный день
func (r *EventRepo) GetDailyEvents(ctx context.Context, date time.Time) ([]entity.Event, error) {
	var events []entity.Event

	day := date.Day()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, userEvents := range r.events {
		for _, event := range userEvents {
			if event.EventTime.Day() == day {
				events = append(events, event)
			}
		}
	}
	return events, nil
}

// GetWeeklyEvents СписокСобытийНаНеделю (дата начала недели);
// Выводит список событий за 7 дней, начиная с дня начала
func (r *EventRepo) GetWeeklyEvents(ctx context.Context, date time.Time) ([]entity.Event, error) {
	var events []entity.Event

	endDay := date.Add(7 * 24 * time.Hour)

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, userEvents := range r.events {
		for _, event := range userEvents {
			if event.EventTime.After(date) && event.EventTime.Before(endDay) {
				events = append(events, event)
			}
		}
	}
	return events, nil
}

// GetMonthlyEvents СписокСобытийНaМесяц (дата начала месяца)
// Выводит список событий за 30 дней, начиная с дня начала
func (r *EventRepo) GetMonthlyEvents(ctx context.Context, date time.Time) ([]entity.Event, error) {
	var events []entity.Event

	endDay := date.Add(7 * 24 * 30 * time.Hour)

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, userEvents := range r.events {
		for _, event := range userEvents {
			if event.EventTime.After(date) && event.EventTime.Before(endDay) {
				events = append(events, event)
			}
		}
	}
	return events, nil
}

// isEventTimeBusy проверка на занятость в заданное время
func (r *EventRepo) isEventTimeBusy(userEvents map[int32]entity.Event, newEvent entity.Event) bool {
	newStartTime := newEvent.EventTime
	newEndTime := newEvent.EventTime.Add(newEvent.Duration)
	for _, userEvent := range userEvents {
		oldStartTime := userEvent.EventTime
		oldEndTime := userEvent.EventTime.Add(userEvent.Duration)
		if (newStartTime.After(oldStartTime) && newStartTime.Before(oldEndTime)) ||
			(newEndTime.After(oldStartTime) && newEndTime.Before(oldEndTime)) {
			return false
		}
	}
	return true
}

// eventValidate проверка ивента на валидность полей
func (r *EventRepo) eventValidate(event entity.Event) error {
	switch {
	case event.Title == "":
		return errapp.ErrEventTitle
	case event.UserId == 0:
		return errapp.ErrEventUserID
	case event.EventTime.IsZero():
		return errapp.ErrEventTime
	case event.Duration == 0:
		return errapp.ErrEventDuration
	}
	return nil
}
