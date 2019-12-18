package timeutil

import "time"

func InRange(rangeStart, rangeEnd, needle time.Time, incStart, incEnd bool) bool {
	rangeStart, rangeEnd = MinMax(rangeStart, rangeEnd)
	if incStart && needle.Before(rangeStart) {
		return false
	} else if needle.Before(rangeStart) || needle.Equal(rangeStart) {
		return false
	}
	if incEnd && needle.After(rangeEnd) {
		return false
	} else if needle.After(rangeEnd) || needle.Equal(rangeEnd) {
		return false
	}
	return true
}
