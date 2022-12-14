package errors

import (
	"errors"
)

var (
	//ErrInvalidEventTitle  = errors.New("invalid event title")
	//ErrInvalidEventTime   = errors.New("invalid event time")
	//ErrInvalidEventUserID = errors.New("invalid event user id")
	ErrEventTimeBusy = errors.New("event time is already busy")
	ErrEventNotFound = errors.New("event not found")
	// ErrInvalidTime        = errors.New("invalid time")
	ErrEventTitle    = errors.New("empty event title")
	ErrEventTime     = errors.New("empty event time")
	ErrEventDuration = errors.New("empty event duration")
)
