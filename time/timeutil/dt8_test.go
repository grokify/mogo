package timeutil

import (
	"testing"
	"time"
)

var dateTime8Tests = []struct {
	dt8    DateTime8
	y      int32
	m      int32
	d      int32
	layout string
	t      string
}{
	{20230613, 2023, 6, 13, DateTextUS, "June 13, 2023"},
	// {20230632, 2023, 6, 32, DateTextUS, "June 32, 2023"},
}

// TestDateTime8 handles dd6 parsing tests
func TestDateTime8(t *testing.T) {
	for _, tt := range dateTime8Tests {
		y, m, d := tt.dt8.Split()
		if y != tt.y || m != tt.m || d != tt.d {
			t.Errorf("timeutil.DateTime8.Split() val (%d) Mismatch: want (%d,%d,%d) got (%d,%d,%d)",
				int(tt.dt8), tt.y, tt.m, tt.d, y, m, d)
		}
		txt, err := tt.dt8.Format(tt.layout)
		if txt != tt.t {
			t.Errorf("timeutil.DateTime8.Format(\"%s\") val (%d) Error: err (%s)",
				tt.layout, int32(tt.dt8), err.Error())
		}
		if txt != tt.t {
			t.Errorf("timeutil.DateTime8.Format(\"%s\") val (%d) Mismatch: want (%s) got (%s)",
				tt.layout, int32(tt.dt8), tt.t, txt)
		}

		dt, err := tt.dt8.Time()
		if err != nil {
			t.Errorf("timeutil.DateTime8.Time() val (%d) Error: err (%s)",
				int32(tt.dt8), err.Error())
		}
		if int(tt.y) != dt.Year() || int(tt.m) != int(dt.Month()) || int(tt.d) != dt.Day() {
			t.Errorf("timeutil.DateTime8.Time() val (%d) Mismatch: want (%d,%d,%d) got (%d,%d,%d)",
				int(tt.dt8), tt.y, tt.m, tt.d, dt.Year(), dt.Month(), dt.Day())
		}
	}
}

var dt8ForStringTests = []struct {
	v    string
	want DateTime8
}{
	{"2006-01-02T15:04:05Z", DateTime8(int32(20060102))}}

func TestDT8ForString(t *testing.T) {
	for _, tt := range dt8ForStringTests {
		got, err := DT8ForString(time.RFC3339, tt.v)
		if err != nil {
			t.Errorf("Dt8ForString(%v): want %v, error %v", tt.v, tt.want, err)
		}
		if got != tt.want {
			t.Errorf("Dt8ForString(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}
