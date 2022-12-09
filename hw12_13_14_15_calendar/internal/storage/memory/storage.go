package memorystorage

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"sync"
	"time"
)

var (
	ErrEventTimeBusy = errors.New("event time is already busy")
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidTime   = errors.New("invalid time")
	ErrEventTitle    = errors.New("empty event title")
	ErrEventTime     = errors.New("empty event time")
	ErrEventDuration = errors.New("empty event duration")
)

type Storage struct {
	events map[int]map[int32]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	m := sync.RWMutex{}
	events := make(map[int]map[int32]storage.Event)
	return &Storage{
		events: events,
		mu:     m,
	}
}

// Create Создать (событие)
func (s *Storage) Create(ctx context.Context, event storage.Event) error {
	if err := s.EventValidate(event); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	event.Id = int32(uuid.New().ID()) // Create unique ID

	userEvents, ok := s.events[event.UserId]
	if !ok {
		s.events[event.UserId] = make(map[int32]storage.Event)
	}

	if !s.IsEventTimeBusy(userEvents, event) {
		return ErrEventTimeBusy
	}

	s.events[event.UserId][event.Id] = event
	return nil
}

// Update Обновить (ID события, событие);
func (s *Storage) Update(ctx context.Context, event storage.Event) error {
	if err := s.EventValidate(event); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	userEvents, ok := s.events[event.UserId]
	if !ok {
		return ErrEventNotFound
	}

	updatedEvent, ok := userEvents[event.Id]
	if !ok {
		return ErrEventNotFound
	}

	if !s.IsEventTimeBusy(userEvents, event) {
		return ErrEventTimeBusy
	}

	updatedEvent.Title = event.Title
	updatedEvent.Desc = event.Desc
	updatedEvent.EventTime = event.EventTime
	updatedEvent.Duration = event.Duration

	s.events[event.UserId][event.Id] = updatedEvent
	return nil
}

// Delete Удалить (ID события);
func (s *Storage) Delete(Id int32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, userEvents := range s.events {
		if _, ok := userEvents[Id]; !ok {
			delete(userEvents, Id)
			return nil
		}
	}
	return ErrEventNotFound
}

// GetDailyEvents СписокСобытийНаДень (дата);
// Выводит все события, которые начинаются в заданный день
func (s *Storage) GetDailyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	var events []storage.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, userEvents := range s.events {
		for _, event := range userEvents {
			if event.EventTime.Day() == date.Day() {
				events = append(events, event)
			}
		}
	}
	return events, nil
}

// GetWeeklyEvents СписокСобытийНаНеделю (дата начала недели);
// Выводит список событий за 7 дней, начиная с дня начала
func (s *Storage) GetWeeklyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	var events []storage.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, userEvents := range s.events {
		for _, event := range userEvents {
			if event.EventTime.Weekday() == date.Weekday() {
				events = append(events, event)
			}
		}
	}
	return events, nil
}

// GetMonthlyEvents СписокСобытийНaМесяц (дата начала месяца)
// Выводит список событий за 30 дней, начиная с дня начала
func (s *Storage) GetMonthlyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	var events []storage.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, userEvents := range s.events {
		for _, event := range userEvents {
			if event.EventTime.Month() == date.Month() {
				events = append(events, event)
			}
		}
	}
	return events, nil
}

// IsEventTimeBusy проверка на занятость в заданное время
func (s *Storage) IsEventTimeBusy(userEvents map[int32]storage.Event, newEvent storage.Event) bool {
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

// EventValidate проверка ивента на валидность полей
func (s *Storage) EventValidate(event storage.Event) error {
	//TODO написать ивент валидатор
	switch {
	case event.Title == "":
		return ErrEventTitle
	case event.EventTime.IsZero():
		return ErrEventTime
	case event.Duration == 0:
		return ErrEventDuration
	}
	return nil
}
