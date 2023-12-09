package timeutil

import (
	"fmt"
	"strings"
	"time"
)

const (
	NanosecondString  = "naonsecond"
	MicrosecondString = "microsecond"
	MillisecondString = "millisecond"
	SecondString      = "second"
	MinuteString      = "minute"
	HourString        = "hour"
	DayString         = "day"
	WeekString        = "week"
	MonthString       = "month"
	QuarterString     = "quarter"
	HalfYearString    = "halfyear"
	YearString        = "year"
	DecadeString      = "decade"
	ScoreString       = "score"
	CenturyString     = "century"
	MillenniaString   = "millennia"
)

type Interval int

const (
	IntervalNanosecond Interval = iota
	IntervalMicrosecond
	IntervalMillisecond
	IntervalSecond
	IntervalMinute
	IntervalHour
	IntervalDay
	IntervalWeek
	IntervalMonth
	IntervalQuarter
	IntervalHalfYear
	IntervalYear
	IntervalDecade
	IntervalScore
	IntervalCentury
	IntervalMillennia
)

var intervals = [...]string{
	NanosecondString,
	MicrosecondString,
	MillisecondString,
	SecondString,
	MinuteString,
	HourString,
	DayString,
	WeekString,
	MonthString,
	QuarterString,
	HalfYearString,
	YearString,
	DecadeString,
	ScoreString,
	CenturyString,
	MillenniaString,
}

func (i Interval) String() string { return intervals[i] }

func ParseInterval(src string) (Interval, error) {
	canonical := strings.ToLower(strings.TrimSpace(src))
	for i, try := range intervals {
		if canonical == try {
			return Interval(i), nil
		}
	}
	return 0, fmt.Errorf("interval [%v] not found", src)
}

func intervalStart(dt time.Time, interval Interval, d time.Weekday) (time.Time, error) {
	switch interval.String() {
	case YearString:
		return yearStart(dt), nil
	case QuarterString:
		return quarterStart(dt), nil
	case MonthString:
		return monthStart(dt), nil
	case WeekString:
		return weekStart(dt, d), nil
	case DayString:
		return dayStart(dt), nil
	default:
		return time.Time{}, fmt.Errorf("interval (%s) not supported in timeutil.IntervalStart", interval.String())
	}
}
