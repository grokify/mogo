package quarter

import (
	"testing"
	"time"
)

var quarterContinuousTests = []struct {
	year     uint64
	quarter  uint64
	quarterc uint64
	rfc3339  string
}{
	{uint64(0), uint64(1), uint64(1), "0000-01-01T00:00:00Z"},
	{uint64(0), uint64(2), uint64(2), "0000-04-01T00:00:00Z"},
	{uint64(0), uint64(3), uint64(3), "0000-07-01T00:00:00Z"},
	{uint64(0), uint64(4), uint64(4), "0000-10-01T00:00:00Z"},
	{uint64(1), uint64(1), uint64(5), "0001-01-01T00:00:00Z"},
	{uint64(1), uint64(4), uint64(8), "0001-10-01T00:00:00Z"},
	{uint64(2), uint64(1), uint64(9), "0002-01-01T00:00:00Z"},
	{uint64(3), uint64(4), uint64(16), "0003-10-01T00:00:00Z"},
	{uint64(4), uint64(1), uint64(17), "0004-01-01T00:00:00Z"},
}

func TestQuarterContinuous(t *testing.T) {
	for _, tt := range quarterContinuousTests {
		gotQuarterC := YearQuarterToQuarterContinuous(tt.year, tt.quarter)
		if gotQuarterC != tt.quarterc {
			t.Errorf("YearQuarterToQuarterContinuous(%v, %v): want [%v], got [%v]",
				tt.year, tt.quarter, tt.quarterc, gotQuarterC)
		}
		dt := QuarterContinuousToTime(tt.quarterc)
		t3339 := dt.Format(time.RFC3339)
		if t3339 != tt.rfc3339 {
			t.Errorf("QuarterContinuousToTime(%v, %v): want [%v], got [%v]",
				tt.year, tt.quarter, tt.rfc3339, t3339)
		}
		wantDt, err := time.Parse(time.RFC3339, tt.rfc3339)
		if err != nil {
			t.Errorf("TimeToQuarterContinuous time.Parse(time.RFC3339, \"%s\") error [%v]",
				tt.rfc3339, err.Error())
		}
		gotMonthc := TimeToQuarterContinuous(wantDt)
		if gotMonthc != tt.quarterc {
			t.Errorf("TimeToQuarterContinuous(\"%s\"): want [%v] got [%v]",
				tt.rfc3339, tt.quarterc, gotQuarterC)
		}
	}
}
