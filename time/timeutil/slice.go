// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

// TimeSlice is used for sorting. e.g.
// sort.Sort(sort.Reverse(timeSlice))
// sort.Sort(timeSlice)
// var times TimeSlice := []time.Time{time.Now()}
type TimeSlice []time.Time

func (s TimeSlice) Less(i, j int) bool { return s[i].Before(s[j]) }
func (s TimeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s TimeSlice) Len() int           { return len(s) }

func Sort(times []time.Time) []time.Time {
	sort.Slice(
		times,
		func(i, j int) bool { return times[i].Before(times[j]) })
	return times
}

func Earliest(times []time.Time, skipZeroes bool) (time.Time, error) {
	if len(times) == 0 {
		return time.Now(), errors.New("No times")
	}
	sort.Slice(
		times,
		func(i, j int) bool { return times[i].Before(times[j]) })

	first := times[0]
	if skipZeroes && TimeIsZeroAny(first) {
		return first, fmt.Errorf("Latest Time is Zero [%s]", first.Format(time.RFC3339))
	}
	return first, nil
}

func Latest(times []time.Time, skipZeroes bool) (time.Time, error) {
	if len(times) == 0 {
		return time.Now(), errors.New("No times")
	}
	sort.Slice(
		times,
		func(i, j int) bool { return times[i].Before(times[j]) })

	last := times[len(times)-1]
	if skipZeroes && TimeIsZeroAny(last) {
		return last, fmt.Errorf("Latest Time is Zero [%s]", last.Format(time.RFC3339))
	}
	return last, nil
}

func QuarterSlice(min, max time.Time) []time.Time {
	min, max = MinMax(min, max)
	minQ := QuarterStart(min)
	maxQ := QuarterStart(max)
	times := []time.Time{}
	cur := minQ
	for cur.Before(maxQ) || cur.Equal(maxQ) {
		times = append(times, cur)
		cur = NextQuarter(cur)
	}
	return times
}
