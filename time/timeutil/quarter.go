package timeutil

import (
	"fmt"
	"math"
	"time"

	"github.com/grokify/mogo/pointer"
)

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

func QuarterAdd(t time.Time, addQuarters int) time.Time {
	if addQuarters == 0 {
		return NewTimeMore(t, 0).QuarterStart()
	}
	y, m := AddMonths(t.Year(), int(t.Month()), addQuarters*3)
	newTime := time.Date(
		y, time.Month(m), 1,
		0, 0, 0, 0,
		pointer.Pointer(pointer.Dereference(t.Location())),
	)
	return NewTimeMore(newTime, 0).QuarterStart()
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
