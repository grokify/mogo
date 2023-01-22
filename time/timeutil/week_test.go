package timeutil

import (
	"testing"
	"time"
)

var weekdayTests = []struct {
	v    int
	want time.Weekday
}{
	{-14, time.Sunday},
	{-13, time.Monday},
	{-1, time.Saturday},
	{0, time.Sunday},
	{6, time.Saturday},
	{7, time.Sunday},
	{15, time.Monday},
	{23, time.Tuesday},
	{29, time.Monday},
}

func TestWeekdayNormalize(t *testing.T) {
	for _, tt := range weekdayTests {
		got := WeekdayNormalized(time.Weekday(tt.v))

		if got != tt.want {
			t.Errorf("mismatch WeekdayNormalize(%d): want [%s], got [%s]", tt.v, tt.want.String(), got.String())
		}
	}
}
