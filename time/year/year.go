package year

import (
	"sort"
	"time"

	"github.com/grokify/mogo/time/timeslice"
	"github.com/grokify/mogo/time/timeutil"
)

// TimeSeriesYear returns time series of months given start and end input times.
func TimeSeriesYear(sortAsc bool, times ...time.Time) timeslice.TimeSlice {
	min, max := timeutil.SliceMinMax(times)
	minYear := timeutil.YearStart(min)
	maxYear := timeutil.YearStart(max)
	timeSeries := timeslice.TimeSlice{}
	curYear := minYear
	for curYear.Before(maxYear) || curYear.Equal(maxYear) {
		timeSeries = append(timeSeries, curYear)
		curYear = curYear.AddDate(1, 0, 0)
	}
	if len(timeSeries) > 1 && sortAsc {
		sort.Sort(timeSeries)
	}
	return timeSeries
}
