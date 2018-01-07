// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"fmt"
	"strconv"
	"time"
)

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

func QuarterInt32Start(yyyyq int32) (time.Time, error) {
	yyyy, q, err := ParseQuarterInt32(yyyyq)
	if err != nil {
		return time.Now(), err
	}
	qm := QuarterToMonth(q)
	t := time.Date(int(yyyy), time.Month(qm), 1, 0, 0, 0, 0, time.UTC)
	return t, nil
}

func ParseQuarterStartString(yyyyqStr string) (time.Time, error) {
	yyyyq, err := strconv.Atoi(yyyyqStr)
	if err != nil {
		return time.Now(), err
	}
	return QuarterInt32Start(int32(yyyyq))
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
	t := time.Date(int(yyyy), time.Month(qm), 0, 23, 59, 59, 0, time.UTC)
	return t, nil
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
