package timeutil

import (
	"strings"
	"testing"
	"time"
)

var sliceMinMaxTests = []struct {
	v   string
	min string
	max string
}{
	{"2006-01-01T00:00:00Z,2006-02-01T00:00:00Z,2006-03-01T00:00:00Z",
		"2006-01-01T00:00:00Z",
		"2006-03-01T00:00:00Z"},
}

// TestDMYHM2ParseTime ensures timeutil.DateDMYHM2 is parsed to GMT timezone.
func TestSliceMinMax(t *testing.T) {
	for _, tt := range sliceMinMaxTests {
		times, err := ParseSlice(strings.Split(tt.v, ","), time.RFC3339)
		if err != nil {
			t.Errorf("time.ParseSlice(%v) Error: [%v]", tt.v, err.Error())
		}
		min, max := SliceMinMax(times)
		if min.Format(time.RFC3339) != tt.min ||
			max.Format(time.RFC3339) != tt.max {
			t.Errorf("timeutil.SliceMinMax(\"%v\") Mismatch: want [%v,%v], got [%v,%v]", tt.v,
				tt.min, tt.max,
				min.Format(time.RFC3339), max.Format(time.RFC3339))
		}
	}
}
