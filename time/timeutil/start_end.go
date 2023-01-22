package timeutil

import (
	"time"
)

func dayStart(dt time.Time) time.Time {
	return time.Date(
		dt.Year(), dt.Month(), dt.Day(),
		0, 0, 0, 0, dt.Location())
}

// weekStart takes a `time.Time` struct and a week start day in the `time.Weekday` format.
func weekStart(t time.Time, d time.Weekday) time.Time {
	ws, err := TimeDeltaDowInt(t, WeekdayNormalized(d), -1, true, true)
	if err != nil {
		panic(err)
	}
	return ws
}

func monthStart(dt time.Time) time.Time {
	return time.Date(
		dt.Year(), dt.Month(), 1,
		0, 0, 0, 0, dt.Location())
}

// quarterStart returns a time.Time for the start of the quarter.
func quarterStart(t time.Time) time.Time {
	qm := QuarterToMonth(MonthToQuarter(uint8(t.Month())))
	return time.Date(t.Year(), time.Month(qm), 1, 0, 0, 0, 0, t.Location())
}

func yearStart(t time.Time) time.Time {
	return time.Date(
		t.Year(), time.January, 1,
		0, 0, 0, 0, t.Location())
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

func dayEnd(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(NanosPerSecondSub1), t.Location())
}

func weekEnd(t time.Time, d time.Weekday) time.Time {
	ws := weekStart(t, d)
	we := ws.AddDate(0, 0, 6)
	return time.Date(we.Year(), we.Month(), we.Day(), 23, 59, 59, int(NanosPerSecondSub1), t.Location())
}

// monthEnd returns a time.Time for the end of the month by second.
func monthEnd(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), MonthEndDay(t.Year(), t.Month()), 23, 59, 59, int(NanosPerSecondSub1), t.Location())
}

// quarterEnd returns a time.Time for the end of the quarter by second.
func quarterEnd(t time.Time) time.Time {
	qs := quarterStart(t)
	qn := TimeDt6AddNMonths(qs, 3)
	return time.Date(qn.Year(), qn.Month(), MonthEndDay(qn.Year(), qn.Month()), 23, 59, 59, int(NanosPerSecondSub1), t.Location())
}

// yearEnd returns a a time.Time for the end of the year.
func yearEnd(t time.Time) time.Time {
	return time.Date(t.Year(), time.December, 31, 23, 59, 59, int(NanosPerSecondSub1), t.Location())
}
