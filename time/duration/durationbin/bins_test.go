package durationbin

import (
	"testing"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

var binFromDurationTests = []struct {
	v           time.Duration
	binDuration time.Duration
}{
	{0 * timeutil.Day, 7 * timeutil.Day},
	{6 * timeutil.Day, 7 * timeutil.Day},
	{7 * timeutil.Day, 7 * timeutil.Day},
	{7*timeutil.Day + 1, 14 * timeutil.Day},
	{1*timeutil.Year - 1, 1 * timeutil.Year},
	{1*timeutil.Year + 0, 1 * timeutil.Year},
	{1*timeutil.Year + 1, 2 * timeutil.Year},
	{1*timeutil.Year + 180*timeutil.Day, 2 * timeutil.Year},
	{2*timeutil.Year - 1, 2 * timeutil.Year},
	{2*timeutil.Year + 0, 2 * timeutil.Year},
	{2*timeutil.Year + 1, 3 * timeutil.Year},
	{3*timeutil.Year + 0, 3 * timeutil.Year},
	{3*timeutil.Year + 2, 4 * timeutil.Year},
	{5*timeutil.Year + 0, 5 * timeutil.Year},
	{5*timeutil.Year + 1, 6 * timeutil.Year},
}

// TestBinForDuration ensures the right bin is returned.
func TestBinFromDuration(t *testing.T) {
	for _, tt := range binFromDurationTests {
		try := BinFromDuration(tt.v)
		if try.Duration != tt.binDuration {
			t.Errorf("durationbins.BinForDuration(%d) Mismatch: want (%d) got (%d)", tt.v,
				tt.binDuration, try.Duration)
		}
	}
}
