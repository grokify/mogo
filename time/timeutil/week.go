package timeutil

import (
	"time"

	"github.com/grokify/mogo/math/mathutil"
)

// FirstDayOfISOWeek returns a time.Time object for the first day of
// a given ISO week.
// https://xferion.com/golang-reverse-isoweek-get-the-date-of-the-first-day-of-iso-week/
func FirstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()

	// iterate back to Monday
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the first week
	for isoYear < year {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the given week
	for isoWeek < week {
		date = date.AddDate(0, 0, 7)
		//isoYear, isoWeek = date.ISOWeek()
	}

	return date
}

func (tm TimeMore) WeekdayNext(d time.Weekday) time.Time {
	today := tm.time.Weekday()
	if d == today {
		return tm.time.Add(NewDuration(7, 0, 0, 0, 0))
	} else if d > today {
		return tm.time.Add(NewDurationFloat(float64(int(d)-int(today)), 0, 0, 0, 0))
	}
	return tm.time.Add(NewDurationFloat(float64(int(today)-int(d)+7), 0, 0, 0, 0))
}

// WeekdayNormalized ensures a `time.Weekday` value is within `[0,6]`. It supports
// converting postiive and negative integers.
func WeekdayNormalized(d time.Weekday) time.Weekday {
	return time.Weekday(mathutil.ModPyInt(int(d), 7))
}
