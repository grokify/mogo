// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"fmt"
	"strconv"
	"time"
)

func InQuarter(dt time.Time, yyyyq int32) (bool, error) {
	qtrStart, err := QuarterInt32StartTime(yyyyq)
	if err != nil {
		return false, err
	}
	return (IsGreaterThan(dt, qtrStart, true) &&
		IsLessThan(dt, QuarterEnd(qtrStart), true)), nil
}

func MustInQuarter(dt time.Time, yyyyq int32) bool {
	qtrStart, err := QuarterInt32StartTime(yyyyq)
	if err != nil {
		panic(err)
	}
	return (IsGreaterThan(dt, qtrStart, true) &&
		IsLessThan(dt, QuarterEnd(qtrStart), true))
}

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
