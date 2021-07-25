package timeslice

import (
	"errors"
	"sort"
	"time"
)

// TimeSlice is used for sorting. e.g.
// sort.Sort(sort.Reverse(timeSlice))
// sort.Sort(timeSlice)
type TimeSlice []time.Time

func (ts TimeSlice) Len() int           { return len(ts) }
func (ts TimeSlice) Less(i, j int) bool { return ts[i].Before(ts[j]) }
func (ts TimeSlice) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }
func (ts TimeSlice) Sort()              { sort.Sort(ts) }

func (ts TimeSlice) Dedupe() TimeSlice {
	newTimeSlice := TimeSlice{}
	seenMap := map[time.Time]int{}
	for _, dt := range ts {
		if _, ok := seenMap[dt]; !ok {
			newTimeSlice = append(newTimeSlice, dt)
			seenMap[dt] = 1
		}
	}
	return newTimeSlice
}

func (ts TimeSlice) Equal(compare TimeSlice) bool {
	if ts.Len() != compare.Len() {
		return false
	}
	for i, dt := range compare {
		if !dt.Equal(ts[i]) {
			return false
		}
	}
	return true
}

func (ts TimeSlice) Format(format string) []string {
	formatted := []string{}
	for _, dt := range ts {
		formatted = append(formatted, dt.Format(format))
	}
	return formatted
}

func (ts TimeSlice) Duplicate() TimeSlice {
	newTs := TimeSlice{}
	for _, t := range ts {
		newTs = append(newTs, t)
	}
	return newTs
}

var (
	EmptyTimeSliceError   = errors.New("empty time slice")
	OutOfBoundsError      = errors.New("out of bounds")
	OutOfBoundsLowerError = errors.New("out of bounds lower")
	OutOfBoundsUpperError = errors.New("out of bounds upper")
)

// RangeLower returns the TimeSlice time value for the range
// lower than or equal to the supplied time.
func (ts TimeSlice) RangeLower(t time.Time, inclusive bool) (time.Time, error) {
	if len(ts) == 0 {
		return t, EmptyTimeSliceError
	}
	if ts[0].After(t) { // lower bound too high
		return t, OutOfBoundsLowerError
	} else if len(ts) == 1 {
		return ts[0], nil
	}
	for i := 1; i < len(ts); i++ {
		if t.Before(ts[i]) {
			return ts[i-1], nil
		} else if inclusive && t.Equal(ts[i]) {
			return ts[i-1], nil
		}
	}
	return ts[len(ts)-1], nil
}

// RangeUpper returns the TimeSlice time value for the range
// lower than or equal to the supplied time. The time `t` must
// be less than or equal to the upper range.
func (ts TimeSlice) RangeUpper(t time.Time, inclusive bool) (time.Time, error) {
	if len(ts) == 0 {
		return t, EmptyTimeSliceError
	} else if ts[len(ts)-1].Before(t) {
		return t, OutOfBoundsUpperError
	}
	for i := range ts {
		tMax := ts[len(ts)-1-i]
		if t.Before(tMax) {
			return tMax, nil
		} else if inclusive && t.Equal(tMax) {
			return tMax, nil
		}
	}
	return ts[0], nil
}

func ParseTimeSlice(format string, times []string) (TimeSlice, error) {
	ts := TimeSlice{}
	for _, raw := range times {
		dt, err := time.Parse(format, raw)
		if err != nil {
			return ts, err
		}
		ts = append(ts, dt)
	}
	return ts, nil
}
