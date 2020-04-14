package timeutil

import (
	"fmt"
	"math"
	"time"
)

// QuarterStart returns a time.Time for the beginning of the
// quarter in UTC time.
func QuarterStart(dt time.Time) time.Time {
	dt = dt.UTC()
	qm := QuarterToMonth(MonthToQuarter(uint8(dt.Month())))
	return time.Date(dt.Year(), time.Month(qm), 1, 0, 0, 0, 0, time.UTC)
}

func QuarterStartString(dt time.Time) string {
	dtStart := QuarterStart(dt)
	return fmt.Sprintf("%v Q%v", dtStart.Year(), MonthToQuarter(uint8(dtStart.Month())))
}

func NextQuarter(dt time.Time) time.Time {
	return TimeDt6AddNMonths(QuarterStart(dt), 3)
}

func DeltaQuarters(dt time.Time, num int) time.Time {
	dt = QuarterStart(dt)
	if num > 0 {
		dt = NextQuarters(dt, uint(num))
	} else if num < 0 {
		dt = PrevQuarters(dt, uint(-1*num))
	}
	return dt
}

func NextQuarters(dt time.Time, num uint) time.Time {
	dt = QuarterStart(dt)
	for i := 0; i < int(num); i++ {
		dt = NextQuarter(dt)
	}
	return dt
}

func PrevQuarter(dt time.Time) time.Time {
	return TimeDt6SubNMonths(QuarterStart(dt), 3)
}

func PrevQuarters(dt time.Time, num uint) time.Time {
	for i := 0; i < int(num); i++ {
		dt = PrevQuarter(dt)
	}
	return dt
}

func IsQuarterStart(t time.Time) bool {
	t = t.UTC()
	if t.Nanosecond() == 0 &&
		t.Second() == 0 &&
		t.Minute() == 0 &&
		t.Hour() == 0 &&
		t.Day() == 1 &&
		(t.Month() == time.January ||
			t.Month() == time.April ||
			t.Month() == time.July ||
			t.Month() == time.October) {
		return true
	}
	return false
}

// MonthToQuarter converts a month to a calendar quarter.
func MonthToQuarter(month uint8) uint8 {
	return uint8(math.Ceil(float64(month) / 3))
}

// QuarterToMonth converts a calendar quarter to a month.
func QuarterToMonth(quarter uint8) uint8 {
	return quarter*3 - 2
}
