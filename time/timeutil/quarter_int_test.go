package timeutil

import (
	"testing"
	"time"
)

var mustInQuarterRangeTests = []struct {
	v     string
	start int32
	end   int32
	in    bool
}{
	{"2019-11-11T00:00:00Z", 20191, 20194, true},
	{"2019-11-11T00:00:00Z", 20191, 20193, false},
}

func TestMustInQuarterRange(t *testing.T) {
	for _, tt := range mustInQuarterRangeTests {
		dt, err := time.Parse(time.RFC3339, tt.v)
		if err != nil {
			t.Errorf("DayTestMustInQuarterRange Date Parse Error [%v]", err.Error())
		}
		in := MustInQuarterRange(dt, tt.start, tt.end)
		if tt.in != in {
			t.Errorf("MustInQuarterRange(\"%v\", %v, %v): want [%v], got [%v]", tt.v, tt.start, tt.end, tt.in, in)
		}
	}
}

var numQuartersInt32Tests = []struct {
	start int32
	end   int32
	num   int
}{
	{20191, 20191, 1},
	{20191, 20194, 4},
	{20191, 20202, 6},
	{20193, 20202, 4},
	{20194, 20201, 2},
	{20193, 20222, 12},
}

func TestNumQuartersInt32(t *testing.T) {
	for _, tt := range numQuartersInt32Tests {
		got, err := NumQuartersInt32(tt.start, tt.end)
		if err != nil {
			t.Errorf("NumQuarters Error [%v]", err.Error())
		}
		if got != tt.num {
			t.Errorf("NumQuarters(%v, %v): want [%v], got [%v]", tt.start, tt.end, tt.num, got)
		}
	}
}

func TestDeltaQuarterInt32(t *testing.T) {
	for _, tt := range numQuartersInt32Tests {
		got, err := QuarterInt32Add(tt.start, tt.num-1)
		if err != nil {
			t.Errorf("DeltaQuarterInt32 Error [%v]", err.Error())
		}
		if got != tt.end {
			t.Errorf("DeltaQuarterInt32(%v, %v): want [%v], got [%v]", tt.start, tt.num, tt.end, got)
		}
	}
}
