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
