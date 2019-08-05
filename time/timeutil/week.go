package timeutil

import (
	"time"
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

func NextWeekday(dow time.Weekday) time.Time {
	now := time.Now()
	today := now.Weekday()
	if dow == today {
		//dur := tu.MustParseDuration("7d")
		return now.Add(NewDurationDays(7))
	} else if dow > today {
		days := int(dow) - int(today)
		//dur := tu.MustParseDuration(fmt.Sprintf("%dd", days))
		return now.Add(NewDurationDays(uint16(days)))
	}
	days := int(today) - int(dow) + 7
	//dur := tu.MustParseDuration(fmt.Sprintf("%dd", days))
	return now.Add(NewDurationDays(uint16(days)))
}
