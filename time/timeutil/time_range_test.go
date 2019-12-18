package timeutil

import (
	"testing"
	"time"
)

var inRangeTests = []struct {
	v       string
	begin   string
	end     string
	inclBeg bool
	inclEnd bool
	inRange bool
}{
	{"2019-12-31T00:00:00Z", "2019-10-01T00:00:00Z", "2021-01-01T00:00:00Z", true, true, true},
	{"2018-12-31T00:00:00Z", "2019-10-01T00:00:00Z", "2021-01-01T00:00:00Z", true, true, false},
}

func TestInRange(t *testing.T) {
	for _, tt := range inRangeTests {
		dtThs, err := time.Parse(time.RFC3339, tt.v)
		if err != nil {
			panic(err)
		}
		dtBeg, err := time.Parse(time.RFC3339, tt.begin)
		if err != nil {
			panic(err)
		}
		dtEnd, err := time.Parse(time.RFC3339, tt.end)
		if err != nil {
			panic(err)
		}
		got := InRange(dtBeg, dtEnd, dtThs, tt.inclBeg, tt.inclEnd)
		if got != tt.inRange {
			t.Errorf("InRange(%v,...): want [%v], got [%v]", tt.v, tt.inRange, got)
		}
	}
}
