package timeutil

import (
	"fmt"
	"time"
)

// TimeSeriesSlice builds a time series based on supplied interval.
func TimeSeriesSlice(interval Interval, times []time.Time) []time.Time {
	if len(times) == 0 {
		return times
	}
	min, max := SliceMinMax(times)
	return TimeSeriesMinMax(interval, min, max)
}

// TimeSeriesMinMax builds a time series based on supplied interval.
func TimeSeriesMinMax(interval Interval, min, max time.Time) []time.Time {
	min, max = MinMax(min, max)
	series := []time.Time{}
	tmMin := NewTimeMore(min, 0)
	tmMax := NewTimeMore(max, 0)
	switch interval {
	case IntervalYear:
		min = tmMin.YearStart()
		max = tmMax.YearStart()
		series = append(series, min)
		cur := min
		for {
			cur = cur.AddDate(1, 0, 0)
			if cur.After(max) {
				break
			}
			series = append(series, cur)
		}
	case IntervalMonth:
		min = tmMin.MonthStart()
		max = tmMin.MonthStart()
		series = append(series, min)
		cur := min
		for {
			cur = cur.AddDate(0, 1, 0)
			if cur.After(max) {
				break
			}
			series = append(series, cur)
		}
	default:
		panic(fmt.Sprintf("interval not supportedin timeutil.TimeSeriesMinMax [%v]", interval))
	}
	return series
}
