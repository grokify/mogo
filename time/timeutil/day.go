package timeutil

import (
	"fmt"
	"time"

	"github.com/grokify/gotilla/type/stringsutil"
)

var days = [...]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

func ParseWeekday(s string) (time.Weekday, error) {
	for i, day := range days {
		if stringsutil.Equal(s, day, true, true) {
			return time.Weekday(i), nil
		}
	}
	return time.Weekday(0), fmt.Errorf("Cannot parse weekday: %s", s)
}

// TimeToDay sets a time's hour, minute, second,
// and nanosecond to 0.
func TimeToDay(dt time.Time) time.Time {
	return time.Date(
		dt.Year(), dt.Month(), dt.Day(),
		0, 0, 0, 0, dt.Location())
}
