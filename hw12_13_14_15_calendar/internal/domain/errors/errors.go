package models

import (
	"errors"
)

var (
	ErrInvalidEventTitle  = errors.New("invalid event title")
	ErrInvalidEventTime   = errors.New("invalid event time")
	ErrInvalidEventUserID = errors.New("invalid event user id")
)
