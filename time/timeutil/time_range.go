package timeutil

import "time"

func InRange(this, start, end time.Time, incStart, incEnd bool) bool {
	start, end = MinMax(start, end)
	if incStart && this.Before(start) {
		return false
	} else if this.Before(start) || this.Equal(start) {
		return false
	}
	if incEnd && this.After(end) {
		return false
	} else if this.After(start) || this.Equal(start) {
		return false
	}
	return true
}
