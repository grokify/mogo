package timeutil

import (
	"testing"
	"time"
)

var monthEndDayTests = []struct {
	year  int
	month time.Month
	day   int
}{
	{2010, time.January, 31},
	{2020, time.January, 31},
	{2023, time.January, 31},
	{2023, time.February, 28},
	{2023, time.June, 30},
	{2020, time.February, 29},
	{2024, time.February, 29},
	{2040, time.February, 29},
}

// TestMonthEndDay returns the last day of a given month and year.
func TestMonthEndDay(t *testing.T) {
	for _, tt := range monthEndDayTests {
		got := MonthEndDay(tt.year, tt.month)
		if got != tt.day {
			t.Errorf("timeutil.MonthEndDay(%d, %d) Mismatch: want (%d), got (%d)", tt.year, tt.month, tt.day, got)
		}
	}
}
