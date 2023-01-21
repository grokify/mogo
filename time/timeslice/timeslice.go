package timeslice

import (
	"errors"
	"sort"
	"time"
)

// Times is used for sorting. e.g.
// sort.Sort(sort.Reverse(timeSlice))
// sort.Sort(timeSlice)
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
	formatted := []string{}
	for _, dt := range ts {
		formatted = append(formatted, dt.Format(format))
	}
	return formatted
}

func (ts Times) Duplicate() Times {
	newTS := Times{}
	newTS = append(newTS, ts...)
	return newTS
}

var (
	ErrEmptyTimeSlice   = errors.New("empty time slice")
	ErrOutOfBounds      = errors.New("out of bounds")
	ErrOutOfBoundsLower = errors.New("out of bounds lower")
	ErrOutOfBoundsUpper = errors.New("out of bounds upper")
)

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
		}
		curRangeLower = nextRangeLower
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
		}
		curRangeUpper = nextRangeUpper
	}
	return sortedTS[0], nil
}

func ParseTimes(format string, times []string) (Times, error) {
	ts := Times{}
	for _, raw := range times {
		dt, err := time.Parse(format, raw)
		if err != nil {
			return ts, err
		}
		ts = append(ts, dt)
	}
	return ts, nil
}
