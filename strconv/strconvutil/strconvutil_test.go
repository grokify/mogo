package strconvutil

import (
	"testing"

	"github.com/grokify/mogo/math/mathutil"
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
		try := mathutil.RoundMore(ChangeToXoXPct(tt.v), 0.5, 0.0)
		// try := mathutil.Round(ChangeToXoXPct(tt.v))
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
		try := mathutil.RoundMore(ChangeToFunnelPct(tt.v), 0.5, 0.0)
		if try != tt.want {
			t.Errorf("strconvutil.ChangeToFunnelPct() Error: with [%v], want [%v], got [%v]",
				tt.v, tt.want, try)
		}
	}
}

var intAbbrTests = []struct {
	v    int64
	want string
}{
	{0, "0"},
	{1, "1"},
	{100, "100"},
	{999, "999"},
	{1000, "1.0K"},
	{1500, "1.5K"},
	{15000, "15K"},
	{150000, "150K"},
	{1200000, "1.2M"},
	{2000000, "2.0M"},
	{20000000, "20M"},
	{200000000, "200M"},
	{2000000000, "2.0B"},
	{20000000000, "20B"},
	{200000000000, "200B"},
	{2500000000000, "2.5T"},
	{25000000000000, "25T"},
	{250000000000000, "250T"},
}

func TestIntAbbrevations(t *testing.T) {
	for _, tt := range intAbbrTests {
		try := Int64Abbreviation(tt.v)
		if try != tt.want {
			t.Errorf("strconvutil.Int64Abbreviation(%v) Error: want [%v], got [%v]",
				tt.v, tt.want, try)
		}
	}
}
