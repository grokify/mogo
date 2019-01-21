package mathutil

import (
	"testing"
)

var ratioIntTests = []struct {
	x1 int
	y1 int
	x2 int
	y2 int
	xo int
	yo int
}{
	{600, 400, 1950, 1300, 1950, 1300},
	{600, 400, 1950, 0, 1950, 1300},
	{600, 400, 0, 1300, 1950, 1300},
	{600, 400, 0, 0, 600, 400},
}

func TestRatioInt(t *testing.T) {
	for _, tt := range ratioIntTests {
		x, y := RatioInt(tt.x1, tt.y1, tt.x2, tt.y2)
		if x != tt.xo || y != tt.yo {
			t.Errorf("mathutil.Ratio Error: Ratio(%v, %v, %v, %v) Want(%v, %v) Got(%v, %v)",
				tt.x1, tt.y1, tt.x2, tt.y2,
				tt.xo, tt.yo,
				x, y)
		}
	}
}
