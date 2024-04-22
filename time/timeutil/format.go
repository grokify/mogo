package timeutil

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// More predefined layouts for use in Time.Format and time.Parse.
const (
	DT14               = "20060102150405"
	DT8                = "20060102"
	DT6                = "200601"
	MonthDay           = "1/2"
	DIN5008FullDate    = "02.01.2006" // German DIN 5008 standard
	DIN5008Date        = "02.01.06"
	RFC3339FullDate    = "2006-01-02"
	RFC3339Milli       = "2006-01-02T15:04:05.999Z07:00"
	RFC3339Dash        = "2006-01-02T15-04-05Z07-00"
	ISO8601            = "2006-01-02T15:04:05Z0700"
	ISO8601TZHour      = "2006-01-02T15:04:05Z07"
	ISO8601NoTZ        = "2006-01-02T15:04:05"
	ISO8601MilliNoTZ   = "2006-01-02T15:04:05.999"
	ISO8601Milli       = "2006-01-02T15:04:05.999Z0700"
	ISO8601CompactZ    = "20060102T150405Z0700"
	ISO8601CompactNoTZ = "20060102T150405"
	ISO8601YM          = "2006-01"
	ISO9075            = "2006-01-02 15:04:05"           // ISO/IEC 9075, MySQL, BigQuery, etc.
	ISO9075MicroTZ     = "2006-01-02 15:04:05.999999-07" // ISO 90875 used by PostgreSQL
	// SQLTimestamp     = "2006-01-02 15:04:05"       // ISO/IEC 9075, MySQL, BigQuery, etc.
	// SQLTimestampPgTz = "2006-01-02 15:04:05.999999-07" // ISO 90875 used by PostgreSQL
	Ruby              = "2006-01-02 15:04:05 -0700" // Ruby Time.now.to_s
	InsightlyAPIQuery = "_1/_2/2006 _3:04:05 PM"
	DateMDY           = "1/2/2006" // an underscore results in a space.
	DateMDYSlash      = "01/02/2006"
	DateDMYDash       = "_2-01-2006"     // Jira XML Date format
	DateDMYHM2        = "02:01:06 15:04" // GMT time in format dd:mm:yy hh:mm
	DateYMD           = RFC3339FullDate
	DateTextUS        = "January _2, 2006"
	DateTextUSAbbr3   = "Jan _2, 2006"
	DateTextEU        = "_2 January 2006"
	DateTextEUAbbr3   = "_2 Jan 2006"
)

const (
	RFC3339Min         = "0000-01-01T00:00:00Z"
	RFC3339Max         = "9999-12-31T23:59:59Z"
	RFC3339Zero        = "0001-01-01T00:00:00Z" // Golang zero value
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
	return "", errors.New("no match")
}

func ParseFirstValueOrZero(layout string, values []string) time.Time {
	dt, err := ParseFirstValue(layout, values)
	if err != nil {
		return TimeZeroRFC3339()
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
		return time.Now(), errors.New("no time values supplied")
	}
	return time.Now(), fmt.Errorf("no valid string of [%v] supplied values", strconv.Itoa(numVals))
}

// ParseOrZero returns a parsed time.Time or the RFC-3339 zero time.
func ParseOrZero(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		return TimeZeroRFC3339()
	}
	return t
}

// ParseFirst attempts to parse a string with a set of layouts.
func ParseFirst(layouts []string, value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if len(value) == 0 || len(layouts) == 0 {
		return time.Now(), fmt.Errorf(
			"requires value [%v] and at least one layout [%v]", value, strings.Join(layouts, ","))
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
	return time.Now(), fmt.Errorf("cannot parse time [%v] with layouts [%v]",
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
		return "", fmt.Errorf("format Not Found: %v", format)
	}
	return format, nil
}

// FormatQuarter takes quarter time and formats it using "Q# YYYY".
func FormatQuarter(t time.Time) string {
	return fmt.Sprintf("Q%d %d", MonthToQuarter(t.Month()), t.Year())
}

// FormatQuarter takes quarter time and formats it using "YYYY Q#".
func FormatQuarterYYYYQ(t time.Time) string {
	return fmt.Sprintf("%d Q%d", t.Year(), MonthToQuarter(t.Month()))
}

// FormatQuarterYYQ takes quarter time and formats it using "'YY Q#".
func FormatQuarterYYQ(t time.Time) string {
	return fmt.Sprintf("%s Q%d", t.Format("'06"), MonthToQuarter(t.Month()))
}

func TimeMinRFC3339() time.Time {
	t0, _ := time.Parse(time.RFC3339, RFC3339Min)
	return t0
}

func TimeZeroRFC3339() time.Time {
	t0, _ := time.Parse(time.RFC3339, RFC3339Zero)
	return t0
}

func TimeZeroUnix() time.Time {
	t0, _ := time.Parse(time.RFC3339, RFC3339ZeroUnix)
	return t0
}

// func (tm TimeMore) IsZero() bool { return tm.Time.Equal(TimeZeroRFC3339()) }
// func IsZeroAny(u time.Time) bool { return TimeIsZeroAny(u) }

func (tm TimeMore) IsZeroAny() bool {
	return isZeroAny(tm.Time())
}

func isZeroAny(u time.Time) bool {
	return u.Equal(TimeZeroRFC3339()) ||
		u.Equal(TimeMinRFC3339()) ||
		u.Equal(TimeZeroUnix())
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
func FormatTimeMulti(dateFormat string, v any) string {
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

func FormatTimeToString(format string) func(time.Time) string {
	return func(dt time.Time) string {
		return dt.Format(format)
	}
}

// OffsetFormat converts an integer offset value to a string
// value to use in string time formats. Note: RFC-3339 times
// use colons and the UTC "Z" offset.
func OffsetFormat(offset int, useColon, useZ bool) string {
	offsetStr := "+0000"
	if offset == 0 {
		if useZ {
			offsetStr = "Z"
		} else if useColon {
			offsetStr = "+00:00"
		}
	} else if offset > 0 {
		if useColon {
			hr := offset / 100
			mn := offset - (hr * 100)
			offsetStr = "+" + fmt.Sprintf("%02d:%02d", hr, mn)
		} else {
			offsetStr = "+" + fmt.Sprintf("%04d", offset)
		}
	} else if offset < 0 {
		if useColon {
			offsetPositive := -1 * offset
			hr := offsetPositive / 100
			mn := offsetPositive - (hr * 100)
			offsetStr = "-" + fmt.Sprintf("%02d:%02d", hr, mn)
		} else {
			offsetStr = "-" + fmt.Sprintf("%04d", -1*offset)
		}
	}
	return offsetStr
}

// var rxSQLTimestamp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}$`)

func ParseTimeUsingOffset(format, raw, sep string, offset int, useColon, useZ bool) (time.Time, error) {
	return time.Parse(format, raw+sep+OffsetFormat(offset, useColon, useZ))
}

// ParseTimeSQLTimestampUsingOffset converts a SQL timestamp without timezone
// adding in a manual integer timezone.
func ParseTimeSQLTimestampUsingOffset(timeStr string, offset int) (time.Time, error) {
	return ParseTimeUsingOffset(Ruby, timeStr, " ", offset, false, false)
	/*
		timeStr = strings.TrimSpace(timeStr)
		if !rxSQLTimestamp.MatchString(timeStr) {
			return time.Now(), fmt.Errorf("E_INVALID_SQL_TIMESTAMP [%v]", timeStr)
		}
		offsetStr := OffsetFormat(offset, useColon, useZ)
		timeStr += " " + offsetStr
		dt, err := time.Parse(Ruby, timeStr)
		return dt, err
	*/
}

const (
	LayoutNameDT4  = "dt4"
	LayoutNameDT6  = "dt6"
	LayoutNameDT8  = "dt8"
	LayoutNameDT14 = "dt14"
)

// IsDTX returns the dtx format if conformant to various DTX values (dt4, dt6, dt8, dt14).
func IsDTX(d int32) (string, error) {
	switch len(strconv.Itoa(int(d))) {
	case 4:
		return LayoutNameDT4, nil
	case 6:
		m := d - ((d / 100) * 100)
		if m < 1 || m > 12 {
			return LayoutNameDT6, errors.New("dt6 length month value is out of bounds")
		}
		return LayoutNameDT6, nil
	case 8:
		dy := d - ((d / 100) * 100)
		if dy < 1 || dy > 31 {
			return LayoutNameDT6, errors.New("dt8 day value is out of bounds")
		}
		return LayoutNameDT8, nil
	case 14:
		return LayoutNameDT14, nil
	default:
		return "", errors.New("length mismatch")
	}
}
