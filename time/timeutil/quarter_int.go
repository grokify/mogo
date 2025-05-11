package timeutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/type/ordered"
)

func InYearQuarter(dt time.Time, yyyyq int) (bool, error) {
	thsQtrStart, err := YearQuarterStartTime(yyyyq)
	if err != nil {
		return false, err
	}
	nxtQtrStart := QuarterAdd(thsQtrStart, 1)
	return (thsQtrStart.Before(dt) || thsQtrStart.Equal(dt)) &&
		nxtQtrStart.After(dt), nil
}

func MustInYearQuarter(dt time.Time, yyyyq int) bool {
	thsQtrStart, err := YearQuarterStartTime(yyyyq)
	if err != nil {
		panic(err)
	}
	nxtQtrStart := QuarterAdd(thsQtrStart, 1)
	return (thsQtrStart.Before(dt) || thsQtrStart.Equal(dt)) &&
		nxtQtrStart.After(dt)
}

// InQuarterRange checks to see if the date is within 2 quarters.
func InQuarterRange(dt time.Time, yyyyq1, yyyyq2 int) (bool, error) {
	dtQ1, err := YearQuarterStartTime(yyyyq1)
	if err != nil {
		return false, err
	}
	dtQ2, err := YearQuarterStartTime(yyyyq2)
	if err != nil {
		return false, err
	}
	dtQ2Next := QuarterAdd(dtQ2, 1)
	dtQ1, _ = MinMax(dtQ1, dtQ2)
	return (dt.Equal(dtQ1) || dt.After(dtQ1)) && (dt.Before(dtQ2Next)), nil
}

// MustInYearQuarterRange returns whether a date is within 2 quarters.
// It panics if the quarter integer is not valid.
func MustInYearQuarterRange(dt time.Time, yyyyq1, yyyyq2 int) bool {
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
	return IsGreaterThan(dt, quarterStart(qtr), true) &&
		IsLessThan(dt, quarterEnd(qtr), true)
}

func EqualQuarter(dt1, dt2 time.Time) bool {
	return YearQuarterForTime(dt1) == YearQuarterForTime(dt2)
}

func YearQuarterForTime(dt time.Time) int {
	dt = dt.UTC()
	q := MonthToQuarter(dt.Month())
	return (dt.Year() * 10) + int(q)
}

func QuarterNow() int { return YearQuarterForTime(time.Now()) }

func ParseQuarterInt32StartEndTimes(yyyyq int) (time.Time, time.Time, error) {
	start, err := YearQuarterStartTime(yyyyq)
	if err != nil {
		return start, start, err
	}
	return start, quarterEnd(start), nil
}

func ParseYearQuarter(yyyyq int) (int, int, error) {
	yyyy := int(float32(yyyyq) / 10.0)
	q := yyyyq - 10*yyyy
	if q < 0 {
		q = -1 * q
	}
	if q < 1 || q > 4 {
		return 0, 0, fmt.Errorf("quarter '%v' is not valid", q)
	}
	return yyyy, q, nil
}

func YearQuarterStartTimeString(yyyyqStr string) (time.Time, error) {
	yyyyq, err := strconv.Atoi(yyyyqStr)
	if err != nil {
		return time.Now(), err
	}
	return YearQuarterStartTime(yyyyq)
}

func YearQuarterEndTimeString(yyyyqStr string) (time.Time, error) {
	yyyyq, err := strconv.Atoi(yyyyqStr)
	if err != nil {
		return time.Now(), err
	}
	return YearQuarterEndTime(yyyyq)
}

func YearQuarterStartTime(yyyyq int) (time.Time, error) {
	yyyy, q, err := ParseYearQuarter(yyyyq)
	if err != nil {
		return time.Now(), err
	}
	qm := QuarterToMonth(Yearquarter(q))
	return time.Date(yyyy, qm, 1, 0, 0, 0, 0, time.UTC), nil
}

func YearQuarterEndTime(yyyyq int) (time.Time, error) {
	qtrBeg, err := YearQuarterStartTime(yyyyq)
	if err != nil {
		return qtrBeg, err
	}
	return quarterEnd(qtrBeg), nil
}

func ParseQuarterStringStartTime(yyyyqStr string) (time.Time, error) {
	yyyyq32, err := strconv.Atoi(yyyyqStr)
	if err != nil {
		return time.Now(), err
	}
	return YearQuarterStartTime(yyyyq32)
}

func QuarterInt32End(yyyyq int) (time.Time, error) {
	yyyy, q, err := ParseYearQuarter(yyyyq)
	if err != nil {
		return time.Now(), err
	}
	if 1 == 4 {
		q = 1
	} else {
		q += 1
	}
	qm := QuarterToMonth(Yearquarter(q))
	return time.Date(yyyy, qm, 0, 23, 59, 59, 0, time.UTC), nil
}

func ParseYearHalf(yyyyh int) (int, int, error) {
	yyyy := int(float32(yyyyh) / 10.0)
	h := yyyyh - 10*yyyy
	if h < 0 {
		h = -1 * h
	}
	if h < 1 || h > 2 {
		return 0, 0, fmt.Errorf("half '%v' is not valid", h)
	}
	return yyyy, h, nil
}

func YearQuarterToYear(yyyyq int) int { return int(float32(yyyyq) / 10) }

func quarterInt32NextSingle(yyyyq int) (int, error) {
	t, err := YearQuarterStartTime(yyyyq)
	if err != nil {
		return 0, err
	}
	tNext := QuarterAdd(t, 1)
	return YearQuarterForTime(tNext), nil
}

func YearQuarterAdd(yyyyq int, numQuarters int) (int, error) {
	if numQuarters < 0 {
		return -1, fmt.Errorf("use positive number of quarters [%v]", numQuarters)
	} else if numQuarters == 0 {
		return yyyyq, nil
	}
	future := yyyyq
	var err error
	for i := 0; i < numQuarters; i++ {
		future, err = quarterInt32NextSingle(future)
		if err != nil {
			return -1, errorsutil.Wrap(err, "future quarter")
		}
	}
	return future, nil
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
	yyyyqNow := YearQuarterForTime(time.Now().UTC())
	if yyyyqSrc == yyyyqNow {
		return time.Now().UTC()
	}
	// return quarter end time
	dtQtrEnd, err := YearQuarterEndTime(yyyyqSrc)
	if err != nil {
		return time.Now().UTC()
	}
	return dtQtrEnd.UTC()
}

var rxYYYYQ = regexp.MustCompile(`^[0-9]{4}[1-4]$`)

func IsYearQuarter(yyyyq int) bool {
	return rxYYYYQ.MatchString(strconv.Itoa(yyyyq))
}

func NumQuartersInt32(start, end int) (int, error) {
	start, end = ordered.MinMax(start, end)
	if !IsYearQuarter(start) {
		return -1, fmt.Errorf("quarterInt32 is not valid [%v] Must end in [1-4]", start)
	}
	if !IsYearQuarter(end) {
		return -1, fmt.Errorf("quarterInt32 is not valid [%v] Must end in [1-4]", end)
	}

	cur := start
	if start == end {
		return 1, nil
	}
	numQuarters := 1
	var err error
	for {
		numQuarters++
		cur, err = YearQuarterAdd(cur, 1)
		if err != nil {
			return -1, err
		}
		if cur == end {
			break
		}
	}
	return numQuarters, nil
}

// QuartersRelToAbs is useful relative date queries.
func QuartersRelToAbs(start, end int) (int, int) {
	if start < 100 {
		start = YearQuarterForTime(
			QuarterAdd(time.Now(), start))
	}
	if end < 100 {
		startTime, err := YearQuarterStartTime(start)
		if err != nil {
			panic(errorsutil.Wrap(err, "timeutil.YearQuarterStartTime"))
		}
		end = YearQuarterForTime(
			QuarterAdd(startTime, end-1))
	}
	return start, end
}
