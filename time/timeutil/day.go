package timeutil

import "time"

// TimeToDay sets a time's hour, minute, second,
// and nanosecond to 0.
func TimeToDay(dt time.Time) time.Time {
	return time.Date(
		dt.Year(), dt.Month(), dt.Day(),
		0, 0, 0, 0, dt.Location())
}
