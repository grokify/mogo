package timeutil

import (
	"fmt"
	"time"
)

// SliceTimeSeries builds a time series based on supplied interval.
func SliceTimeSeries(times []time.Time, interval Interval) []time.Time {
	if len(times) == 0 {
		return times
	}
	min, max := SliceMinMax(times)
	return MinMaxTimeSeries(min, max, interval)
}

// MinMaxTimeSeries builds a time series based on supplied interval.
func MinMaxTimeSeries(min, max time.Time, interval Interval) []time.Time {
	min, max = MinMax(min, max)
	series := []time.Time{}
	switch interval {
	case Month:
		min = MonthStart(min)
		max = MonthStart(max)
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
		panic(fmt.Sprintf("E_INTERVAL_NOT_SUPPORTED [%v]", interval))
	}
	return series
}
