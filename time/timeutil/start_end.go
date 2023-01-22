package timeutil

import (
	"time"
)

func dayStart(dt time.Time) time.Time {
	return time.Date(
		dt.Year(), dt.Month(), dt.Day(),
		0, 0, 0, 0, dt.Location())
}

// weekStart takes a time.Time object and a week start day in the time.Weekday format.
func weekStart(dt time.Time, dow time.Weekday) (time.Time, error) {
	return TimeDeltaDowInt(dt.UTC(), int(dow), -1, true, true)
}

func monthStart(dt time.Time) time.Time {
	return time.Date(
		dt.Year(), dt.Month(), 1,
		0, 0, 0, 0, dt.Location())
}

// quarterStart returns a time.Time for the beginning of the quarter in UTC time.
func quarterStart(dt time.Time) time.Time {
	dt = dt.UTC()
	qm := QuarterToMonth(MonthToQuarter(uint8(dt.Month())))
	return time.Date(dt.Year(), time.Month(qm), 1, 0, 0, 0, 0, time.UTC)
}

func yearStart(dt time.Time) time.Time {
	return time.Date(
		dt.Year(), time.January, 1,
		0, 0, 0, 0, dt.Location())
}

/*
func NextYearStart(dt time.Time) time.Time {
	return time.Date(
		dt.Year()+1, time.January, 1,
		0, 0, 0, 0, dt.Location())
}
*/

func isMonthStart(dt time.Time) bool {
	return dt.Day() == 1 &&
		dt.Hour() == 0 &&
		dt.Minute() == 0 &&
		dt.Second() == 0 &&
		dt.Nanosecond() == 0
}

func isQuarterStart(t time.Time) bool {
	t = t.UTC()
	return t.Nanosecond() == 0 &&
		t.Second() == 0 &&
		t.Minute() == 0 &&
		t.Hour() == 0 &&
		t.Day() == 1 &&
		(t.Month() == time.January ||
			t.Month() == time.April ||
			t.Month() == time.July ||
			t.Month() == time.October)
}

func isYearStart(dt time.Time) bool {
	return dt.Month() == time.January && isMonthStart(dt)
}

// quarterEnd returns a time.Time for the end of the quarter by second in UTC time.
func quarterEnd(dt time.Time) time.Time {
	qs := quarterStart(dt.UTC())
	qn := TimeDt6AddNMonths(qs, 3)
	return time.Date(qn.Year(), qn.Month(), 0, 23, 59, 59, 999999999, time.UTC)
}

// yearEnd returns a a time.Time for the end of the year in UTC time.
func yearEnd(dt time.Time) time.Time {
	return time.Date(dt.UTC().Year(), time.December, 31, 23, 59, 59, 999999999, time.UTC)
}
