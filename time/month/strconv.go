package month

import (
	"errors"

	"github.com/grokify/mogo/sort/sortutil"
	"github.com/grokify/mogo/time/timeutil"
)

// DT6FormatOrDefault coverts an `int32` to a `layout`, using `def`ault. Panics on
// `time.Parse()` error.
func MustDT6FormatOrDefault(layout string, dt6 int32, def string) string {
	str, err := DT6FormatOrDefault(layout, dt6, def)
	if err != nil {
		panic(err)
	}
	return str
}

// DT6FormatOrDefault coverts an `int32` to a `layout`, using `def`ault
func DT6FormatOrDefault(layout string, dt6 int32, def string) (string, error) {
	if dt6 <= 0 {
		return def, nil
	}
	dt, err := timeutil.TimeForDT6(dt6)
	if err != nil {
		return "", err
	}
	return dt.Format(layout), nil
}

// StartEndDT6s returns a start and end from a range. `-1` is returned if a 0 lenth slice
// is provided. If only a start date is provided, a `0` is returned to represent "not provided",
// which can indicate something like "the present".
func StartEndDT6s(dt6s []int32) (int32, int32) {
	if len(dt6s) == 0 {
		return -1, -1
	}

	// dt6s = slicesutil.Dedupe(dt6s) // Support if start and end months are the identical.
	sortutil.Slice(dt6s)

	if dt6s[0] <= 0 {
		if len(dt6s) > 1 {
			return dt6s[1], 0
		} else {
			return -1, -1
		}
	}

	if len(dt6s) > 1 {
		return dt6s[0], dt6s[len(dt6s)-1]
	} else {
		return dt6s[0], 0
	}
}

// DT6sToString returns a formatted date range. This function is can be used to generate a string such as:
// ` (Jan 2007-Dec2008)` with the call `DT6sToString(200701, 200812, "Jan 2006", "Present", " (", "-", ")")` or
// ` (Jan 2007-Present)` with the call `DT6sToString(200701, 0, "Jan 2006", "Present", " (", "-", ")")`.
func DT6sToString(startDT6, endDT6 int32, layout, presentText, prefixText, joinText, suffixText string) (string, error) {
	if startDT6 <= 0 && endDT6 <= 0 {
		return "", nil
	}
	startText, endText, err := DT6sToStrings(startDT6, endDT6, layout, presentText)
	if err != nil {
		return "", err
	}
	if startText != "" && endText != "" {
		return prefixText + startText + joinText + endText + suffixText, nil
	} else if startText != "" {
		return prefixText + startText + suffixText, nil
	} else if endText != "" {
		return prefixText + endText + suffixText, nil
	}
	return "", nil
}

func DT6sToStrings(startDT6, endDT6 int32, layout, presentText string) (string, string, error) {
	if startDT6 <= 0 && endDT6 <= 0 {
		return "", "", nil
	}
	if startDT6 <= 0 && endDT6 > 0 {
		return "", "", errors.New("start date cannot be undefined (lte zero)")
	}
	startDT6Fmt, err := DT6FormatOrDefault(layout, startDT6, presentText)
	if err != nil {
		return startDT6Fmt, "", err
	}
	endDT6Fmt, err := DT6FormatOrDefault(layout, endDT6, presentText)
	return startDT6Fmt, endDT6Fmt, err
}
