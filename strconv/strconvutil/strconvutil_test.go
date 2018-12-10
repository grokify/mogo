package strconvutil

import (
	"testing"

	"github.com/grokify/gotilla/math/mathutil"
)

var changeToXoXPctTests = []struct {
	v    float64
	want float64
}{
	{1.15, 15.0},
	{1.1, 10.0},
	{1.0, 0.0},
	{0.9, -10.0},
	{0.87, -13.0},
}

func TestChangeToXoXPctTests(t *testing.T) {
	for _, tt := range changeToXoXPctTests {
		// without math.Round, we end up with:
		// Error: with [0.9], want [-10], got [-9.999999999999998]
		try := mathutil.Round(ChangeToXoXPct(tt.v), 0.5, 0.0)
		if try != tt.want {
			t.Errorf("strconvutil.ChangeToXoXPct() Error: with [%v], want [%v], got [%v]",
				tt.v, tt.want, try)
		}
	}
}

var changeToFunnelPctTests = []struct {
	v    float64
	want float64
}{
	{2.0, 200.0},
	{1.5, 150.0},
	{1.15, 115.0},
	{1.1, 110.0},
	{1.0, 100.0},
	{0.9, 90.0},
	{0.87, 87.0},
	{0.5, 50.0},
	{0.25, 25.0},
}

func TestChangeToFunnelPctTests(t *testing.T) {
	for _, tt := range changeToFunnelPctTests {
		// without math.Round, we end up with:
		// Error: with [0.9], want [-10], got [-9.999999999999998]
		try := mathutil.Round(ChangeToFunnelPct(tt.v), 0.5, 0.0)
		if try != tt.want {
			t.Errorf("strconvutil.ChangeToFunnelPct() Error: with [%v], want [%v], got [%v]",
				tt.v, tt.want, try)
		}
	}
}
