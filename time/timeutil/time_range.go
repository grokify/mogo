package timeutil

import "time"

// InRange checks to see if a time is within a provided time range
// with options whether the begin and end ranges are inclusive or
// exclusive. Exclusive ranges are the default.
func InRange(rangeStart, rangeEnd, needle time.Time, incStart, incEnd bool) bool {
	rangeStart, rangeEnd = MinMax(rangeStart, rangeEnd)
	if (incStart && needle.Before(rangeStart)) ||
		(needle.Before(rangeStart) || needle.Equal(rangeStart)) ||
		(incEnd && needle.After(rangeEnd)) ||
		(needle.After(rangeEnd) || needle.Equal(rangeEnd)) {
		return false
	}
	return true
}
