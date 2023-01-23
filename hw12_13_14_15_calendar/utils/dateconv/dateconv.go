package dateconv

import (
	"fmt"
	"time"
)

const (
	_apiDayLayout  = "2006-01-02"
	_apiTimeLayout = "2006-01-02 15:04:05"
)

func TimeToString(date time.Time) string {
	return date.Format(_apiTimeLayout)
}

func StringToTime(s string) (time.Time, error) {
	date, err := time.Parse(_apiTimeLayout, s)
	if err != nil {
		return date, fmt.Errorf("time.Parse(_apiDateLayout, s): %w", err)
	}
	return date, nil
}

func StringToDay(s string) (time.Time, error) {
	date, err := time.Parse(_apiDayLayout, s)
	if err != nil {
		return date, fmt.Errorf("time.Parse(_apiDayLayout, s): %w", err)
	}
	return date, nil
}
