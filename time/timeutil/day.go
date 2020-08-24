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
