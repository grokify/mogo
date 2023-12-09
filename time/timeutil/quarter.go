package timeutil

import (
	"fmt"
	"math"
	"time"
)

/*
// QuarterStart returns a time.Time for the beginning of the
// quarter in UTC time.
func QuarterStart(dt time.Time) time.Time {
	dt = dt.UTC()
	qm := QuarterToMonth(MonthToQuarter(uint8(dt.Month())))
	return time.Date(dt.Year(), time.Month(qm), 1, 0, 0, 0, 0, time.UTC)
}
*/

func QuarterStartString(t time.Time) string {
	dtStart := NewTimeMore(t, 0).QuarterStart()
	return fmt.Sprintf("%v Q%v", dtStart.Year(), MonthToQuarter(dtStart.Month()))
}

/*
func DeltaQuarters(dt time.Time, num int) time.Time {
	dt = QuarterStart(dt)
	if num > 0 {
		dt = QuarterNext(dt, uint(num))
	} else if num < 0 {
		dt = QuarterPrev(dt, uint(-1*num))
	}
	return dt
}
*/

func quarterNextSingle(t time.Time) time.Time {
	return TimeDT6AddNMonths(NewTimeMore(t, 0).QuarterStart(), 3)
}

func QuarterAdd(t time.Time, count int) time.Time {
	if count == 0 {
		return NewTimeMore(t, 0).QuarterStart()
	} else if count < 0 {
		return quarterPrev(t, uint(-1*count))
	}
	return quarterNext(t, uint(count))
}

func quarterNext(t time.Time, count uint) time.Time {
	t = NewTimeMore(t, 0).QuarterStart()
	for i := 0; i < int(count); i++ {
		t = quarterNextSingle(t)
	}
	return t
}

func quarterPrevSingle(t time.Time) time.Time {
	return TimeDT6AddNMonths(NewTimeMore(t, 0).QuarterStart(), -3)
}

func quarterPrev(t time.Time, num uint) time.Time {
	for i := 0; i < int(num); i++ {
		t = quarterPrevSingle(t)
	}
	return t
}

// MonthToQuarter converts a month to a calendar quarter.
func MonthToQuarter(month time.Month) Yearquarter {
	// func MonthToQuarter(month uint8) uint8 {
	// return uint8(math.Ceil(float64(month) / 3))
	return Yearquarter(uint8(math.Ceil(float64(month) / 3)))
}

// QuarterToMonth converts a calendar quarter to a month.
func QuarterToMonth(quarter Yearquarter) time.Month {
	// func QuarterToMonth(quarter uint8) uint8 {
	//	return quarter*3 - 2
	return time.Month(int(quarter)*3 - 2)
}
