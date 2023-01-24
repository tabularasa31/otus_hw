package memoryrepo

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	errapp "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
)

// EventRepo -.
type EventRepo struct {
	events map[int]map[int]entity.EventDB
	mu     *sync.RWMutex
}

// New -.
func New() *EventRepo {
	m := sync.RWMutex{}
	events := make(map[int]map[int]entity.EventDB)
	return &EventRepo{
		events: events,
		mu:     &m,
	}
}

// CreateEvent Создать (событие).
func (r *EventRepo) CreateEvent(_ context.Context, eventDB *entity.EventDB) (*entity.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create unique event ID
	eventDB.ID = int(uuid.New().ID())

	_, ok := r.events[eventDB.UserID]
	if !ok {
		r.events[eventDB.UserID] = make(map[int]entity.EventDB)
	}
	r.events[eventDB.UserID][eventDB.ID] = *eventDB
	res := r.events[eventDB.UserID][eventDB.ID]

	return res.Dto(), nil
}

// UpdateEvent Обновить (событие).
func (r *EventRepo) UpdateEvent(_ context.Context, eventDB *entity.EventDB) (*entity.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	userEvents, ok := r.events[eventDB.UserID]
	if !ok {
		return nil, errapp.ErrEventNotFound
	}

	updatedEvent, ok := userEvents[eventDB.ID]
	if !ok {
		return nil, errapp.ErrEventNotFound
	}

	updatedEvent.Title = eventDB.Title
	updatedEvent.Desc = eventDB.Desc
	updatedEvent.StartTime = eventDB.StartTime
	updatedEvent.EndTime = eventDB.EndTime

	r.events[eventDB.UserID][eventDB.ID] = updatedEvent
	res := r.events[eventDB.UserID][eventDB.ID]
	return res.Dto(), nil
}

// DeleteEvent Удалить (ID события).
func (r *EventRepo) DeleteEvent(_ context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, userEvents := range r.events {
		if _, ok := userEvents[id]; !ok {
			delete(userEvents, id)
			return nil
		}
	}
	return errapp.ErrEventNotFound
}

// GetEventsByDates СписокСобытийЗаПериод (дата).
// Выводит все события за период, которые начинаются в заданный день.
func (r *EventRepo) GetEventsByDates(_ context.Context, uid int, start time.Time, end time.Time) ([]entity.Event, error) {
	var userEvents []entity.Event

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, userEvent := range r.events[uid] {
		if userEvent.StartTime.After(start) && userEvent.StartTime.Before(end) {
			userEvents = append(userEvents, *userEvent.Dto())
		}
	}

	return userEvents, nil
}

// GetAllEventsByTime СписокСобытийЗаПериод (дата).
// Выводит все события за период, которые начинаются в заданный день.
func (r *EventRepo) GetAllEventsByTime(_ context.Context, start time.Time) ([]entity.Event, error) {
	var result []entity.Event

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, all := range r.events {
		for _, event := range all {
			if event.StartTime == start {
				result = append(result, *event.Dto())
			}
		}
	}

	return result, nil
}

func (r *EventRepo) DeleteOldEventsFromRepo(_ context.Context, date time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, userEvents := range r.events {
		for _, event := range userEvents {
			if event.StartTime.Before(date) {
				delete(r.events, event.ID)
				return nil
			}
		}
	}
	return errapp.ErrEventNotFound
}
