package year

import (
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

// TimesYearStarts returns time series of months given start and end input times.
func TimesYearStarts(times ...time.Time) timeutil.Times {
	timesStarts := timeutil.Times{}
	if len(times) == 0 {
		return timesStarts
	}
	min, max := timeutil.SliceMinMax(times)
	minYear := timeutil.NewTimeMore(min, 0).YearStart()
	maxYear := timeutil.NewTimeMore(max, 0).YearStart()
	curYear := minYear
	for curYear.Before(maxYear) || curYear.Equal(maxYear) {
		timesStarts = append(timesStarts, curYear)
		curYear = curYear.AddDate(1, 0, 0)
	}
	return timesStarts
}
