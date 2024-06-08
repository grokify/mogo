package timeutil

import (
	"errors"
	"sort"
	"time"
)

var (
	ErrEmptyTimeSlice   = errors.New("empty time slice")
	ErrOutOfBounds      = errors.New("out of bounds")
	ErrOutOfBoundsLower = errors.New("out of bounds lower")
	ErrOutOfBoundsUpper = errors.New("out of bounds upper")
)

// Times is a slice of `time.Time` that supports a number of functions and can be used for sorting,
// e.g. `sort.Sort(times)` or `sort.Sort(sort.Reverse(times))`.
type Times []time.Time

func (ts Times) Len() int           { return len(ts) }
func (ts Times) Less(i, j int) bool { return ts[i].Before(ts[j]) }
func (ts Times) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }
func (ts Times) Sort()              { sort.Sort(ts) }
func (ts Times) SortReverse()       { sort.Sort(sort.Reverse(ts)) }

func (ts Times) Dedupe() Times {
	newTimeSlice := Times{}
	seenMap := map[time.Time]int{}
	for _, dt := range ts {
		if _, ok := seenMap[dt]; !ok {
			newTimeSlice = append(newTimeSlice, dt)
			seenMap[dt] = 1
		}
	}
	return newTimeSlice
}

func (ts Times) Deltas() Durations {
	var durs Durations
	if len(ts) < 2 {
		return durs
	}
	for i := 1; i < len(ts); i++ {
		durs = append(durs, ts[i].Sub(ts[i-1]))
	}
	return durs
}

func (ts Times) Duplicate() Times {
	var newTS Times
	newTS = append(newTS, ts...)
	return newTS
}

func (ts Times) Equal(compare Times) bool {
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

func (ts Times) Format(format string) []string {
	var formatted []string
	for _, dt := range ts {
		formatted = append(formatted, dt.Format(format))
	}
	return formatted
}

func (ts Times) IsSorted(asc bool) bool {
	deltas := ts.Deltas()
	if len(deltas) == 0 {
		return true
	}
	if asc {
		for delta := range deltas {
			if delta < 0 {
				return false
			}
		}
		return true
	}
	for delta := range deltas {
		if delta > 0 {
			return false
		}
	}
	return true
}

func (ts Times) Max() *time.Time {
	if len(ts) == 0 {
		return nil
	}
	var max time.Time
	for i, t := range ts {
		if i == 0 {
			max = t
		} else if t.Sub(max) > 0 {
			max = t
		}
	}
	return &max
}

func (ts Times) Min() *time.Time {
	if len(ts) == 0 {
		return nil
	}
	var min time.Time
	for i, t := range ts {
		if i == 0 {
			min = t
		} else if t.Sub(min) < 0 {
			min = t
		}
	}
	return &min
}

// RangeLower returns the TimeSlice time value for the range
// lower than or equal to the supplied time.
func (ts Times) RangeLower(t time.Time, inclusive bool) (time.Time, error) {
	if len(ts) == 0 {
		return t, ErrEmptyTimeSlice
	}
	sortedTS := ts.Dedupe()
	sort.Sort(sortedTS)

	if sortedTS[0].After(t) {
		return t, ErrOutOfBoundsLower
	}

	curRangeLower := sortedTS[0]
	for _, nextRangeLower := range sortedTS {
		if t.Before(nextRangeLower) {
			return curRangeLower, nil
		} else if inclusive && t.Equal(nextRangeLower) {
			return nextRangeLower, nil
		} else {
			curRangeLower = nextRangeLower
		}
	}
	return sortedTS[len(sortedTS)-1], nil
}

// RangeUpper returns the Times time value for the range
// lower than or equal to the supplied time. The time `t` must
// be less than or equal to the upper range.
func (ts Times) RangeUpper(t time.Time, inclusive bool) (time.Time, error) {
	if len(ts) == 0 {
		return t, ErrEmptyTimeSlice
	}
	sortedTS := ts.Dedupe()
	sort.Sort(sortedTS)

	if sortedTS[len(sortedTS)-1].Before(t) {
		return t, ErrOutOfBoundsUpper
	}
	curRangeUpper := sortedTS[len(sortedTS)-1]
	for i := range sortedTS {
		// check times in reverse order
		nextRangeUpper := sortedTS[len(sortedTS)-1-i]
		if t.After(nextRangeUpper) {
			return curRangeUpper, nil
		} else if inclusive && t.Equal(nextRangeUpper) {
			return nextRangeUpper, nil
		} else {
			curRangeUpper = nextRangeUpper
		}
	}
	return sortedTS[0], nil
}

func ParseTimes(format string, times []string) (Times, error) {
	var ts Times
	for _, raw := range times {
		if dt, err := time.Parse(format, raw); err != nil {
			return ts, err
		} else {
			ts = append(ts, dt)
		}
	}
	return ts, nil
}
