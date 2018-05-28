// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"fmt"
	"strings"
	"time"
)

// More predefined layouts for use in Time.Format and time.Parse.
const (
	DT14             = "20060102150405"
	DT8              = "20060102"
	DT6              = "200601"
	RFC3339FullDate  = "2006-01-02"
	ISO8601YM        = "2006-01"
	ISO8601ZHour     = "2006-01-02T15:04:05-07"
	ISO8601Z         = "2006-01-02T15:04:05-0700"
	ISO8601ZCompact  = "20060102T150405Z"
	ISO8601MilliNoTz = "2006-01-02T15:04:05.000"
	DateMDYSlash     = "01/02/2006"
	DateDMYHM2       = "02:01:06 15:04" // GMT time in format dd:mm:yy hh:mm
)

const (
	RFC3339Min         = "0000-01-01T00:00:00Z"
	RFC3339Max         = "9999-12-31T23:59:59Z"
	RFC3339Zero        = "0001-01-01T00:00:00Z"
	RFC3339YMDZeroUnix = int64(-62135596800)
)

// Reformat a time string from one format to another
func FromTo(value, fromLayout, toLayout string) (string, error) {
	t, err := time.Parse(fromLayout, strings.TrimSpace(value))
	if err != nil {
		return "", err
	}
	return t.Format(toLayout), nil
}

// ParseFirst attempts to parse a string with a set of layouts.
func ParseFirst(layouts []string, value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if len(value) == 0 || len(layouts) == 0 {
		return time.Now(), fmt.Errorf(
			"Requires value [%v] and at least one layout [%v]", value, strings.Join(layouts, ","))
	}
	for _, layout := range layouts {
		layout = strings.TrimSpace(layout)
		if len(layout) == 0 {
			continue
		}
		if dt, err := time.Parse(layout, value); err == nil {
			return dt, nil
		}
	}
	return time.Now(), fmt.Errorf("Cannot parse time [%v] with layouts [%v]",
		value, strings.Join(layouts, ","))
}

var FormatMap = map[string]string{
	"RFC3339":    time.RFC3339,
	"RFC3339YMD": RFC3339FullDate,
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
	tt, isNil, err := timeUnmarshalJSON(buf, RFC3339FullDate)
	if err != nil || isNil {
		return err
	}
	t.Time = tt
	return nil
}

func (t RFC3339YMDTime) MarshalJSON() ([]byte, error) {
	return timeMarshalJSON(t.Time, RFC3339FullDate)
}

func (t *ISO8601NoTzMilliTime) UnmarshalJSON(buf []byte) error {
	tt, isNil, err := timeUnmarshalJSON(buf, ISO8601MilliNoTz)
	if err != nil || isNil {
		return err
	}
	t.Time = tt
	return nil
}

func (t ISO8601NoTzMilliTime) MarshalJSON() ([]byte, error) {
	return timeMarshalJSON(t.Time, ISO8601MilliNoTz)
}

func timeUnmarshalJSON(buf []byte, layout string) (time.Time, bool, error) {
	str := string(buf)
	isNil := true
	if str == "null" || str == "\"\"" {
		return time.Time{}, isNil, nil
	}
	tt, err := time.Parse(layout, strings.Trim(str, `"`))
	if err != nil {
		return time.Time{}, false, err
	}
	return tt, false, nil
}

func timeMarshalJSON(t time.Time, layout string) ([]byte, error) {
	return []byte(`"` + t.Format(layout) + `"`), nil
}
