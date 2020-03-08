// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// More predefined layouts for use in Time.Format and time.Parse.
const (
	DT14               = "20060102150405"
	DT8                = "20060102"
	DT6                = "200601"
	MonthDay           = "1/2"
	RFC3339FullDate    = "2006-01-02"
	ISO8601            = "2006-01-02T15:04:05Z0700"
	ISO8601TZHour      = "2006-01-02T15:04:05Z07"
	ISO8601NoTZ        = "2006-01-02T15:04:05"
	ISO8601MilliNoTZ   = "2006-01-02T15:04:05.000"
	ISO8601CompactZ    = "20060102T150405Z0700"
	ISO8601CompactNoTZ = "20060102T150405"
	ISO8601YM          = "2006-01"
	InsightlyApiQuery  = "_1/_2/2006 _3:04:05 PM"
	SQLTimestamp       = "2006-01-02 15:04:05" // MySQL, BigQuery, etc.
	DateMDY            = "1/_2/2006"
	DateMDYSlash       = "01/02/2006"
	DateDMYHM2         = "02:01:06 15:04" // GMT time in format dd:mm:yy hh:mm
)

const (
	RFC3339Min         = "0000-01-01T00:00:00Z"
	RFC3339Max         = "9999-12-31T23:59:59Z"
	RFC3339Zero        = "0001-01-01T00:00:00Z"
	RFC3339ZeroUnix    = "1970-01-01T00:00:00Z"
	RFC3339YMDZeroUnix = int64(-62135596800)
)

// Reformat a time string from one format to another
// Deprecated...
func FromTo(value, fromLayout, toLayout string) (string, error) {
	t, err := time.Parse(fromLayout, strings.TrimSpace(value))
	if err != nil {
		return "", err
	}
	return t.Format(toLayout), nil
}

// Reformat a time string from one format to another
func FromTo2(fromLayout, toLayout, value string) (string, error) {
	t, err := time.Parse(fromLayout, strings.TrimSpace(value))
	if err != nil {
		return "", err
	}
	return t.Format(toLayout), nil
}

func FromToFirstValueOrEmpty(fromLayout, toLayout string, values []string) string {
	dtString, err := FromToFirstValue(fromLayout, toLayout, values)
	if err != nil {
		return ""
	}
	return dtString
}

func FromToFirstValue(fromLayout, toLayout string, values []string) (string, error) {
	for _, val := range values {
		dt, err := time.Parse(fromLayout, val)
		if err == nil {
			return dt.Format(toLayout), nil
		}
	}
	return "", errors.New("No Match")
}

func ParseFirstValueOrZero(layout string, values []string) time.Time {
	dt, err := ParseFirstValue(layout, values)
	if err != nil {
		return TimeRFC3339Zero()
	}
	return dt
}

func ParseFirstValue(layout string, values []string) (time.Time, error) {
	for _, val := range values {
		try, err := time.Parse(layout, val)
		if err == nil {
			return try, nil
		}
	}
	numVals := len(values)
	if numVals == 0 {
		return time.Now(), errors.New("No Time Values Supplied")
	}
	return time.Now(), fmt.Errorf("No Valid String of [%v] Supplied Values", strconv.Itoa(numVals))
}

// ParseOrZero returns a parsed time.Time or the RFC-3339 zero time.
func ParseOrZero(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		return TimeRFC3339Zero()
	}
	return t
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

// FormatQuarter takes quarter time and formats it using "Q# YYYY".
func FormatQuarterYYYYQ(t time.Time) string {
	return fmt.Sprintf("%d Q%d", t.Year(), MonthToQuarter(uint8(t.Month())))
}

func TimeRFC3339Min() time.Time {
	t0, _ := time.Parse(time.RFC3339, RFC3339Min)
	return t0
}

func TimeRFC3339Zero() time.Time {
	t0, _ := time.Parse(time.RFC3339, RFC3339Zero)
	return t0
}

func TimeRFC3339ZeroUnix() time.Time {
	t0, _ := time.Parse(time.RFC3339, RFC3339ZeroUnix)
	return t0
}

func IsZero(u time.Time) bool { return TimeIsZeroAny(u) }

func IsZeroAny(u time.Time) bool { return TimeIsZeroAny(u) }

func TimeIsZeroAny(u time.Time) bool {
	if u.Equal(TimeRFC3339Zero()) ||
		u.Equal(TimeRFC3339Min()) ||
		u.Equal(TimeRFC3339ZeroUnix()) {
		return true
	}
	return false
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
	tt, isNil, err := timeUnmarshalJSON(buf, ISO8601MilliNoTZ)
	if err != nil || isNil {
		return err
	}
	t.Time = tt
	return nil
}

func (t ISO8601NoTzMilliTime) MarshalJSON() ([]byte, error) {
	return timeMarshalJSON(t.Time, ISO8601MilliNoTZ)
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

func ParseSlice(layout string, strings []string) ([]time.Time, error) {
	times := []time.Time{}
	for _, raw := range strings {
		t, err := time.Parse(layout, raw)
		if err != nil {
			return times, err
		}
		times = append(times, t)
	}
	return times, nil
}

// FormatTimeMulti formats a `time.Time` object or
// an epoch number. It is adapted from `github.com/wcharczuk/go-chart`.
func FormatTimeMulti(dateFormat string, v interface{}) string {
	if typed, isTyped := v.(time.Time); isTyped {
		return typed.Format(dateFormat)
	}
	if typed, isTyped := v.(int64); isTyped {
		return time.Unix(0, typed).Format(dateFormat)
	}
	if typed, isTyped := v.(float64); isTyped {
		return time.Unix(0, int64(typed)).Format(dateFormat)
	}
	return ""
}
