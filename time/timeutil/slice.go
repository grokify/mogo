package timeutil

import (
	"errors"
	"fmt"
	"sort"
	"strings"
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

func FirstNonZeroTimeParsed(format string, times []string) (time.Time, error) {
	if len(times) == 0 {
		return time.Now(), errors.New("No times provided")
	}
	for _, timeString := range times {
		timeString = strings.TrimSpace(timeString)
		if len(timeString) == 0 {
			continue
		}
		t, err := time.Parse(format, timeString)
		if err != nil {
			return t, err
		}
		if !IsZeroAny(t) {
			return t, nil
		}
	}
	return time.Now(), errors.New("No times provided")
}

func FirstNonZeroTime(times ...time.Time) (time.Time, error) {
	if len(times) == 0 {
		return time.Now(), errors.New("No times provided")
	}
	for _, t := range times {
		if !IsZeroAny(t) {
			return t, nil
		}
	}
	return time.Now(), errors.New("No times provided")
}

func MustFirstNonZeroTime(times ...time.Time) time.Time {
	t, err := FirstNonZeroTime(times...)
	if err != nil {
		return TimeZeroRFC3339()
	}
	return t
}

func TimeSliceMinMax(times []time.Time) (time.Time, time.Time, error) {
	min := TimeZeroRFC3339()
	max := TimeZeroRFC3339()
	if len(times) == 0 {
		return min, max, errors.New("timeutil.TimeSliceMinMax provided with empty slice")
	}
	for i, t := range times {
		if i == 0 {
			min = t
			max = t
		} else {
			if t.Before(min) {
				min = t
			}
			if t.After(max) {
				max = t
			}
		}
	}
	return min, max, nil
}

// Distinct returns a time slice with distinct, or unique, times.
func Distinct(times []time.Time) []time.Time {
	filtered := []time.Time{}
	seen := map[time.Time]int{}
	for _, t := range times {
		if _, ok := seen[t]; !ok {
			seen[t] = 1
			filtered = append(filtered, t)
		}
	}
	return filtered
}
