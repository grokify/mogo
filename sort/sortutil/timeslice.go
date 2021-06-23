package sortutil

import (
	"sort"
	"time"
)

// TimeSlice is a time slice construct that can be used with
// `sort.Sort`.
type TimeSlice []time.Time

func (p TimeSlice) Len() int           { return len(p) }
func (p TimeSlice) Less(i, j int) bool { return p[i].Before(p[j]) }
func (p TimeSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TimeSlice) Sort()              { sort.Sort(p) }

func (p TimeSlice) Dedupe() TimeSlice {
	newTimeSlice := TimeSlice{}
	seenMap := map[time.Time]int{}
	for _, dt := range p {
		if _, ok := seenMap[dt]; !ok {
			newTimeSlice = append(newTimeSlice, dt)
			seenMap[dt] = 1
		}
	}
	return newTimeSlice
}

func (p TimeSlice) Equal(compare TimeSlice) bool {
	if p.Len() != compare.Len() {
		return false
	}
	for i, dt := range compare {
		if !dt.Equal(p[i]) {
			return false
		}
	}
	return true
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

func (p TimeSlice) Format(format string) []string {
	formatted := []string{}
	for _, dt := range p {
		formatted = append(formatted, dt.Format(format))
	}
	return formatted
}
