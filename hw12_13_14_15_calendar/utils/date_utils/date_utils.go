package date_utils

import (
	"fmt"
	"time"
)

const (
	// _dbDateLayout  = time.RFC3339
	_apiDateLayout = "2006-01-02 15:04:05"
)

func TimeToString(date time.Time) string {
	return date.Format(_apiDateLayout)
}

func StringToTime(s string) (time.Time, error) {
	date, err := time.Parse(_apiDateLayout, s)
	if err != nil {
		return date, fmt.Errorf("time.Parse(_apiDateLayout, s): %w", err)
	}
	return date, nil
}

func DurationToString(d time.Duration) string {
	return d.String()
}

func StringToDuration(s string) (time.Duration, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return d, nil
}
