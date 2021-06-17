package timeutil

import (
	"fmt"
	"strings"
	"time"
)

type Interval int

const (
	Decade Interval = iota
	Year
	Quarter
	Month
	Week
	Day
	Hour
	Minute
	Second
	Millisecond
	Microsecond
	Nanosecond
)

var intervals = [...]string{
	"decade",
	"year",
	"quarter",
	"month",
	"week",
	"day",
	"hour",
	"minute",
	"second",
	"millisecond",
	"microsecond",
	"nanosecond",
}

func (i Interval) String() string { return intervals[i] }

func ParseInterval(src string) (Interval, error) {
	canonical := strings.ToLower(strings.TrimSpace(src))
	for i, try := range intervals {
		if canonical == try {
			return Interval(i), nil
		}
	}
	return Year, fmt.Errorf("Interval [%v] not found.", src)
}

func IntervalStart(dt time.Time, interval Interval, dow time.Weekday) (time.Time, error) {
	switch interval.String() {
	case "year":
		return YearStart(dt), nil
	case "quarter":
		return QuarterStart(dt), nil
	case "month":
		return MonthStart(dt), nil
	case "week":
		return WeekStart(dt, dow)
	default:
		return time.Time{}, fmt.Errorf("Interval [%v] not supported in timeutil.IntervalStart.", interval)
	}
}
