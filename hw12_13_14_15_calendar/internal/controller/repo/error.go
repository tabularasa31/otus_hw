package repo

import (
	"errors"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrEventTitle    = errors.New("empty event title")
	ErrEventUserID   = errors.New("empty event user id")
	ErrEventTime     = errors.New("empty event time")
	ErrEventDuration = errors.New("empty event duration")
)
