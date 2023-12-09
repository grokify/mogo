package timeutil

import "time"

const (
	Day         = 24 * time.Hour
	Week        = 7 * Day
	WorkDay     = 8 * time.Hour
	WorkWeek    = 5 * WorkDay
	WorkDay996  = 12 * time.Hour
	WorkWeek996 = 6 * WorkDay996

	HoursPerDay = float32(24)
	DaysPerWeek = float32(7)

	DaySeconds  = 24 * 60 * 60
	WeekSeconds = 7 * DaySeconds
	YearSeconds = (365 * DaySeconds) + (6 * 60 * 60)

	MonthsEN = `["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]`
)

// A Yearquarter specifies a quarter of the year (Winter = 1, ...).
type Yearquarter uint8

const (
	Winter Yearquarter = 1 + iota
	Spring
	Summer
	Autumn
)

/*
// String returns the English name of the quarter ("Winter", "Spring", ...).
func (Yearquarter) String() string {
	if January <= m && m <= December {
		return longMonthNames[m-1]
	}
	buf := make([]byte, 20)
	n := fmtInt(buf, uint64(m))
	return "%!Month(" + string(buf[n:]) + ")"
}
*/
