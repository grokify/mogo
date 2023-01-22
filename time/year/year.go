package year

import (
	"sort"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

// TimeSeriesYear returns time series of months given start and end input times.
func TimeSeriesYear(sortAsc bool, times ...time.Time) timeutil.Times {
	min, max := timeutil.SliceMinMax(times)
	minYear := timeutil.YearStart(min)
	maxYear := timeutil.YearStart(max)
	timeSeries := timeutil.Times{}
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
