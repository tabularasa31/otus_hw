package memorystorage

import (
	"errors"
	"github.com/google/uuid"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"sync"
	"time"
)

const (
	layout = "02-01-2006 15:04" // day-month-year hour:min
	day    = "02-01-2006"       // day-month-year
)

var (
	ErrEventTimeBusy = errors.New("event time is already busy")
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidTime   = errors.New("invalid time")
)

type Storage struct {
	events map[int]map[uuid.UUID]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	m := sync.RWMutex{}
	events := make(map[int]map[uuid.UUID]storage.Event)
	return &Storage{
		events: events,
		mu:     m,
	}
}

// Create Создать (событие);
func (s *Storage) Create(event storage.Event) error {
	if err := event.EventValidate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	event.Id = uuid.New()

	userEvents, ok := s.events[event.UserId]
	if !ok {
		s.events[event.UserId] = make(map[uuid.UUID]storage.Event)
	}

	if !s.IsEventTimeBusy(userEvents, event) {
		return ErrEventTimeBusy
	}

	s.events[event.UserId][event.Id] = event
	return nil
}

// Update Обновить (ID события, событие);
func (s *Storage) Update(event storage.Event) error {
	if err := event.EventValidate(); err != nil {
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
func (s *Storage) Delete(Id uuid.UUID) error {
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
func (s *Storage) GetDailyEvents(date string) ([]storage.Event, error) {
	check, err := time.Parse(layout, date)
	if err != nil {
		return nil, ErrInvalidTime
	}

	start, err := time.Parse(day, date)
	if err != nil {
		return nil, ErrInvalidTime
	}

	end := start.Add(time.Hour*23 + time.Minute*59)

	return s.eventsInTimeSpan(start, end, check), nil
}

// GetWeeklyEvents СписокСобытийНаНеделю (дата начала недели);
// Выводит список событий за 7 дней, начиная с дня начала
func (s *Storage) GetWeeklyEvents(date string) ([]storage.Event, error) {
	check, err := time.Parse(layout, date)
	if err != nil {
		return nil, ErrInvalidTime
	}

	start, err := time.Parse(day, date)
	if err != nil {
		return nil, ErrInvalidTime
	}

	end := start.Add(time.Hour*167 + time.Minute*59)

	return s.eventsInTimeSpan(start, end, check), nil
}

// GetMonthlyEvents СписокСобытийНaМесяц (дата начала месяца)
// Выводит список событий за 30 дней, начиная с дня начала
func (s *Storage) GetMonthlyEvents(date string) ([]storage.Event, error) {
	check, err := time.Parse(layout, date)
	if err != nil {
		return nil, ErrInvalidTime
	}

	start, err := time.Parse(day, date)
	if err != nil {
		return nil, ErrInvalidTime
	}

	end := start.Add(time.Hour*719 + time.Minute*59)

	return s.eventsInTimeSpan(start, end, check), nil
}

func (s *Storage) IsEventTimeBusy(userEvents map[uuid.UUID]storage.Event, newEvent storage.Event) bool {
	for _, userEvent := range userEvents {
		if newEvent.EventTime.After(userEvent.EventTime) && newEvent.EventTime.Before(userEvent.EventTime.Add(userEvent.Duration)) {
			return false
		}
	}
	return true
}

func (s *Storage) eventsInTimeSpan(start, end, check time.Time) []storage.Event {
	var events []storage.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, userEvents := range s.events {
		for _, event := range userEvents {
			if check.After(start) && check.Before(end) {
				events = append(events, event)
			}
		}
	}

	return events
}
