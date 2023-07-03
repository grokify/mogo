package timeutil

import (
	"testing"
)

var durationTests = []struct {
	day  int
	hour int
	min  int
	sec  int
	nsec int
	sum  int64
}{
	{0, 1, 1, 1, 1, 3661000000001},
	{1, 1, 1, 1, 1, 90061000000001},
	{2, 2, 2, 2, 2, 180122000000002},
	{2, 0, 0, 0, 0, NanosPerDay * 2},
	{1, 0, 0, 0, 0, NanosPerDay},
	{0, 1, 0, 0, 0, NanosPerHour},
	{0, 0, 1, 0, 0, NanosPerMinute},
	{0, 0, 0, 1, 0, NanosPerSecond},
	{0, 0, 0, 0, 1, 1},
}

func TestNewDuration(t *testing.T) {
	for _, tt := range durationTests {
		got := NewDuration(tt.day, tt.hour, tt.min, tt.sec, tt.nsec)
		if int64(got) != tt.sum {
			t.Errorf("timeutil.TimeToDd6(%d,%d,%d,%d,%d) Mismatch: want (%d) got (%d)",
				tt.day, tt.hour, tt.min, tt.sec, tt.nsec, tt.sum, got)
		}
	}
}
