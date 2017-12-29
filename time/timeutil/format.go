// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"fmt"
	"strings"
	"time"
)

const (
	DT14               = "20060102150405"
	DT6                = "200601"
	DT8                = "20060102"
	RFC3339Min         = "0000-01-01T00:00:00Z"
	RFC3339Max         = "9999-12-31T23:59:59Z"
	RFC3339Zero        = "0001-01-01T00:00:00Z"
	RFC3339YMD         = "2006-01-02"
	RFC3339YMDZeroUnix = int64(-62135596800)
	ISO8601YM          = "2006-01"
	ISO8601Z2          = "2006-01-02T15:04:05-07"
	ISO8601Z4          = "2006-01-02T15:04:05-0700"
	ISO8601ZCompact    = "20060102T150405Z"
	ISO8601NoTzMilli   = "2006-01-02T15:04:05.000"
)

var FormatMap = map[string]string{
	"RFC3339":    time.RFC3339,
	"RFC3339YMD": RFC3339YMD,
	"ISO8601YM":  ISO8601YM,
}

func GetFormat(formatName string) (string, error) {
	format, ok := FormatMap[strings.TrimSpace(formatName)]
	if !ok {
		return "", fmt.Errorf("Format Not Found: %v", format)
	}
	return format, nil
}

// FormatQuarter takes quarter time and formats it using "Q# YYYY".
func FormatQuarter(t time.Time) string {
	return fmt.Sprintf("Q%d %d", MonthToQuarter(uint8(t.Month())), t.Year())
}

func TimeRFC3339Zero() time.Time {
	t0, _ := time.Parse(time.RFC3339, RFC3339Zero)
	return t0
}

type RFC3339YMDTime struct{ time.Time }

type ISO8601NoTzMilliTime struct{ time.Time }

func (t *RFC3339YMDTime) UnmarshalJSON(buf []byte) error {
	tt, isNil, err := timeUnmarshalJSON(buf, RFC3339YMD)
	if err != nil || isNil {
		return err
	}
	t.Time = tt
	return nil
}

func (t RFC3339YMDTime) MarshalJSON() ([]byte, error) {
	return timeMarshalJSON(t.Time, RFC3339YMD)
}

func (t *ISO8601NoTzMilliTime) UnmarshalJSON(buf []byte) error {
	tt, isNil, err := timeUnmarshalJSON(buf, ISO8601NoTzMilli)
	if err != nil || isNil {
		return err
	}
	t.Time = tt
	return nil
}

func (t ISO8601NoTzMilliTime) MarshalJSON() ([]byte, error) {
	return timeMarshalJSON(t.Time, ISO8601NoTzMilli)
}

func timeUnmarshalJSON(buf []byte, format string) (time.Time, bool, error) {
	str := string(buf)
	isNil := true
	if str == "null" || str == "\"\"" {
		return time.Time{}, isNil, nil
	}
	tt, err := time.Parse(format, strings.Trim(str, `"`))
	if err != nil {
		return time.Time{}, false, err
	}
	return tt, false, nil
}

func timeMarshalJSON(t time.Time, format string) ([]byte, error) {
	return []byte(`"` + t.Format(format) + `"`), nil
}
