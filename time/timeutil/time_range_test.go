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

var parseTimeRangeIntervalTests = []struct {
	v    string
	tMin string
	tMax string
}{
	{"2023Q4", "2023-10-01T00:00:00Z", "2023-12-31T23:59:59.999999999Z"},
	{"2023H2", "2023-07-01T00:00:00Z", "2023-12-31T23:59:59.999999999Z"},
	{"2024H1", "2024-01-01T00:00:00Z", "2024-06-30T23:59:59.999999999Z"},
	{"2024M1", "2024-01-01T00:00:00Z", "2024-01-31T23:59:59.999999999Z"},
}

func TestParseTimeRangeInterval(t *testing.T) {
	for _, tt := range parseTimeRangeIntervalTests {
		tr, err := ParseTimeRangeInterval(tt.v)
		if err != nil {
			t.Errorf("ParseTimeRangeInterval(\"%s\"): error (%s)", tt.v, err.Error())
		}
		tMin := MustParse(time.RFC3339Nano, tt.tMin).UTC()
		if !tr.Min.Equal(tMin) {
			t.Errorf("ParseTimeRangeInterval(\"%s\"): mismatch tr.Min want (%s) got (%s)", tt.v, tt.tMin, tMin.Format(time.RFC3339Nano))
		}
		tMax := MustParse(time.RFC3339Nano, tt.tMax).UTC()
		if !tr.Max.Equal(tMax) {
			t.Errorf("ParseTimeRangeInterval(\"%s\"): mismatch tr.Max want (%s) got (%s)", tt.v, tt.tMax, tMax.Format(time.RFC3339Nano))
		}
	}
}
