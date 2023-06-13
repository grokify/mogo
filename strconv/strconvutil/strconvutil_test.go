package strconvutil

import (
	"testing"

	"github.com/grokify/mogo/math/mathutil"
	"golang.org/x/exp/constraints"
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

type testItoa[T constraints.Integer] struct {
	// see: https://stackoverflow.com/questions/68166558/generic-structs-with-go
	v    T
	want string
}

var itoaInt64Tests = []testItoa[int64]{
	// testItoa[int8]{v: 127, want: "256"},
	{v: -9223372036854775808, want: "-9223372036854775808"},
	{v: 9223372036854775807, want: "9223372036854775807"},
}

var itoaUint64Tests = []testItoa[uint64]{
	// testItoa[int8]{v: 127, want: "256"},
	{0, "0"},
	{18446744073709551615, "18446744073709551615"},
}

func TestItoa(t *testing.T) {
	for _, tt := range itoaInt64Tests {
		try := Itoa(tt.v)
		if try != tt.want {
			t.Errorf("strconvutil.Itoa(%d) Mismatch: want (%v), got (%v)",
				tt.v, tt.want, try)
		}
	}
	for _, tt := range itoaUint64Tests {
		try := Itoa(tt.v)
		if try != tt.want {
			t.Errorf("strconvutil.Itoa(%d) Mismatch: want (%v), got (%v)",
				tt.v, tt.want, try)
		}
	}
}
