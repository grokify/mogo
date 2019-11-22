// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func InQuarter(dt time.Time, yyyyq int32) (bool, error) {
	thsQtrStart, err := QuarterInt32StartTime(yyyyq)
	if err != nil {
		return false, err
	}
	nxtQtrStart := NextQuarter(thsQtrStart)
	return (thsQtrStart.Before(dt) || thsQtrStart.Equal(dt)) &&
		nxtQtrStart.After(dt), nil
}

func MustInQuarter(dt time.Time, yyyyq int32) bool {
	thsQtrStart, err := QuarterInt32StartTime(yyyyq)
	if err != nil {
		panic(err)
	}
	nxtQtrStart := NextQuarter(thsQtrStart)
	return (thsQtrStart.Before(dt) || thsQtrStart.Equal(dt)) &&
		nxtQtrStart.After(dt)
}

// InQuarterRange checks to see if the date is within 2 quarters.
func InQuarterRange(dt time.Time, yyyyq1, yyyyq2 int32) (bool, error) {
	dtQ1, err := QuarterInt32StartTime(yyyyq1)
	if err != nil {
		return false, err
	}
	dtQ2, err := QuarterInt32StartTime(yyyyq2)
	if err != nil {
		return false, err
	}
	dtQ2Next := NextQuarter(dtQ2)
	dtQ1, dtQ2 = MinMax(dtQ1, dtQ2)
	return (dt.Equal(dtQ1) || dt.After(dtQ1)) && (dt.Before(dtQ2Next)), nil
}

// MustInQuarterRange returns whether a date is within 2 quarters.
// It panics if the quarter integer is not valid.
func MustInQuarterRange(dt time.Time, yyyyq1, yyyyq2 int32) bool {
	inRange, err := InQuarterRange(dt, yyyyq1, yyyyq2)
	if err != nil {
		panic(err)
	}
	return inRange
}

/*
func InQuarterOld(dt time.Time, yyyyq int32) (bool, error) {
	qtrStart, err := QuarterInt32StartTime(yyyyq)
	if err != nil {
		return false, err
	}
	return (IsGreaterThan(dt, qtrStart, true) &&
		IsLessThan(dt, NextQuarter(qtrStart), false)), nil
}

func MustInQuarterOld(dt time.Time, yyyyq int32) bool {
	qtrStart, err := QuarterInt32StartTime(yyyyq)
	if err != nil {
		panic(err)
	}
	return (IsGreaterThan(dt, qtrStart, true) &&
		IsLessThan(dt, QuarterEnd(qtrStart), true))
}
*/

func InQuarterTime(dt, qtr time.Time) bool {
	return IsGreaterThan(dt, QuarterStart(qtr), true) &&
		IsLessThan(dt, QuarterEnd(qtr), true)
}

func EqualQuarter(dt1, dt2 time.Time) bool {
	return QuarterInt32ForTime(dt1) == QuarterInt32ForTime(dt2)
}

func QuarterInt32ForTime(dt time.Time) int32 {
	dt = dt.UTC()
	q := MonthToQuarter(uint8(dt.Month()))
	return (int32(dt.Year()) * int32(10)) + int32(q)
}

func QuarterInt32Now() int32 { return QuarterInt32ForTime(time.Now()) }

func PraseQuarterInt32StartEndTimes(yyyyq int32) (time.Time, time.Time, error) {
	start, err := QuarterInt32StartTime(yyyyq)
	if err != nil {
		return start, start, err
	}
	return start, QuarterEnd(start), nil
}

func ParseQuarterInt32(yyyyq int32) (int32, uint8, error) {
	yyyy := int32(float32(yyyyq) / 10.0)
	q := yyyyq - 10*yyyy
	if q < 0 {
		q = -1 * q
	}
	if q < 1 || q > 4 {
		return int32(0), uint8(0), fmt.Errorf("Quarter '%v' is not valid", q)
	}
	return yyyy, uint8(q), nil
}

func QuarterStringStartTime(yyyyqStr string) (time.Time, error) {
	yyyyq, err := strconv.Atoi(yyyyqStr)
	if err != nil {
		return time.Now(), err
	}
	return QuarterInt32StartTime(int32(yyyyq))
}

func QuarterStringEndTime(yyyyqStr string) (time.Time, error) {
	yyyyq, err := strconv.Atoi(yyyyqStr)
	if err != nil {
		return time.Now(), err
	}
	return QuarterInt32EndTime(int32(yyyyq))
}

func QuarterInt32StartTime(yyyyq int32) (time.Time, error) {
	yyyy, q, err := ParseQuarterInt32(yyyyq)
	if err != nil {
		return time.Now(), err
	}
	qm := QuarterToMonth(q)
	return time.Date(int(yyyy), time.Month(qm), 1, 0, 0, 0, 0, time.UTC), nil
}

func QuarterInt32EndTime(yyyyq int32) (time.Time, error) {
	qtrBeg, err := QuarterInt32StartTime(yyyyq)
	if err != nil {
		return qtrBeg, err
	}
	return QuarterEnd(qtrBeg), nil
}

func ParseQuarterStringStartTime(yyyyqStr string) (time.Time, error) {
	yyyyq, err := strconv.Atoi(yyyyqStr)
	if err != nil {
		return time.Now(), err
	}
	return QuarterInt32StartTime(int32(yyyyq))
}

func QuarterInt32End(yyyyq int32) (time.Time, error) {
	yyyy, q, err := ParseQuarterInt32(yyyyq)
	if err != nil {
		return time.Now(), err
	}
	if 1 == 4 {
		q = 1
	} else {
		q += 1
	}
	qm := QuarterToMonth(q)
	return time.Date(int(yyyy), time.Month(qm), 0, 23, 59, 59, 0, time.UTC), nil
}

func ParseHalf(yyyyh int32) (int32, uint8, error) {
	yyyy := int32(float32(yyyyh) / 10.0)
	h := yyyyh - 10*yyyy
	if h < 0 {
		h = -1 * h
	}
	if h < 1 || h > 2 {
		return int32(0), uint8(0), fmt.Errorf("Half '%v' is not valid", h)
	}
	return yyyy, uint8(h), nil
}

func QuarterInt32ToYear(yyyyq int32) int32 { return int32(float32(yyyyq) / 10) }

func NextQuarterInt32(yyyyq int32) (int32, error) {
	t, err := QuarterInt32StartTime(yyyyq)
	if err != nil {
		return int32(0), err
	}
	tNext := NextQuarter(t)
	return QuarterInt32ForTime(tNext), nil
}

// AnyStringToQuarterTime returns the current time if in the
// current quarter or the end of any previous quarter.
func AnyStringToQuarterTime(yyyyqSrcStr string) time.Time {
	yyyyqSrcStr = strings.TrimSpace(yyyyqSrcStr)
	// If not a string, return now time.
	if len(yyyyqSrcStr) != 5 {
		return time.Now().UTC()
	}
	// If not a yyyyq pattern, return now time.
	rx := regexp.MustCompile(`^[0-9]{4}[1-4]$`)
	m := rx.FindString(strings.TrimSpace(yyyyqSrcStr))
	if len(m) != 5 {
		return time.Now().UTC()
	}
	// If cannot parse to integer, return now time.
	yyyyqSrc, err := strconv.Atoi(yyyyqSrcStr)
	if err != nil {
		return time.Now().UTC()
	}
	// Have good yyyyq
	// If yyyySrc == yyyyNow, return now time.
	yyyyqNow := QuarterInt32ForTime(time.Now().UTC())
	if int32(yyyyqSrc) == yyyyqNow {
		return time.Now().UTC()
	}
	// return quarter end time
	dtQtrEnd, err := QuarterInt32EndTime(int32(yyyyqSrc))
	if err != nil {
		return time.Now().UTC()
	}
	return dtQtrEnd.UTC()
}
