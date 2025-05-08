package quarter

import (
	"testing"
	"time"
)

var quarterContinuousTests = []struct {
	year     uint32
	quarter  uint32
	quarterc uint32
	rfc3339  string
}{
	{uint32(0), uint32(1), uint32(1), "0000-01-01T00:00:00Z"},
	{uint32(0), uint32(2), uint32(2), "0000-04-01T00:00:00Z"},
	{uint32(0), uint32(3), uint32(3), "0000-07-01T00:00:00Z"},
	{uint32(0), uint32(4), uint32(4), "0000-10-01T00:00:00Z"},
	{uint32(1), uint32(1), uint32(5), "0001-01-01T00:00:00Z"},
	{uint32(1), uint32(4), uint32(8), "0001-10-01T00:00:00Z"},
	{uint32(2), uint32(1), uint32(9), "0002-01-01T00:00:00Z"},
	{uint32(3), uint32(4), uint32(16), "0003-10-01T00:00:00Z"},
	{uint32(4), uint32(1), uint32(17), "0004-01-01T00:00:00Z"},
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
			t.Errorf("TimeToQuarterContinuous time.Parse(time.RFC3339, \"%s\") error [%s]",
				tt.rfc3339, err.Error())
		}
		gotQuarterC, err = TimeToQuarterContinuous(wantDt)
		if err != nil {
			t.Errorf("TimeToQuarterContinuous(\"%s\") error [%s]",
				wantDt.Format(time.RFC3339), err.Error())
		}
		if gotQuarterC != tt.quarterc {
			t.Errorf("TimeToQuarterContinuous(\"%s\"): want [%d] got [%d]",
				tt.rfc3339, tt.quarterc, gotQuarterC)
		}
	}
}
