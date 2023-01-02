package date_utils

import "time"

const (
	_dbLayout  = time.RFC3339
	_apiLayout = "2006-01-02 15:04:05"
)

func TimeToString(date time.Time) string {
	return date.Format(_apiLayout)
}

func StringToTime(s string) (time.Time, error) {
	date, err := time.Parse(_dbLayout, s)
	if err != nil {
		return time.Time{}, err
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
