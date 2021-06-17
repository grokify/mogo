package timeutil

import (
	"time"
)

func DayStart(dt time.Time) time.Time {
	return time.Date(
		dt.Year(), dt.Month(), dt.Day(),
		0, 0, 0, 0, dt.Location())
}

// WeekStart takes a time.Time object and a week start day
// in the time.Weekday format.
func WeekStart(dt time.Time, dow time.Weekday) (time.Time, error) {
	return TimeDeltaDowInt(dt.UTC(), int(dow), -1, true, true)
}

func MonthStart(dt time.Time) time.Time {
	return time.Date(
		dt.Year(), dt.Month(), 1,
		0, 0, 0, 0, dt.Location())
}

func YearStart(dt time.Time) time.Time {
	return time.Date(
		dt.Year(), time.January, 1,
		0, 0, 0, 0, dt.Location())
}

func NextYearStart(dt time.Time) time.Time {
	return time.Date(
		dt.Year()+1, time.January, 1,
		0, 0, 0, 0, dt.Location())
}

func IsYearStart(dt time.Time) bool {
	if dt.Month() == time.January &&
		dt.Day() == 1 &&
		dt.Hour() == 0 &&
		dt.Minute() == 0 &&
		dt.Second() == 0 &&
		dt.Nanosecond() == 0 {
		return true
	}
	return false
}

// QuarterEnd returns a time.Time for the end of the
// quarter by second in UTC time.
func QuarterEnd(dt time.Time) time.Time {
	qs := QuarterStart(dt.UTC())
	qn := TimeDt6AddNMonths(qs, 3)
	return time.Date(qn.Year(), qn.Month(), 0, 23, 59, 59, 0, time.UTC)
}

// YearEnd returns a a time.Time for the end of the year in UTC time.
func YearEnd(dt time.Time) time.Time {
	return time.Date(dt.UTC().Year(), time.December, 31, 23, 59, 59, 999999999, time.UTC)
}
