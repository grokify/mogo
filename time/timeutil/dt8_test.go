package timeutil

import (
	"encoding/json"
	"testing"
	"time"
)

var dateTime8Tests = []struct {
	dt8    DateTime8
	y      uint32
	m      uint32
	d      uint32
	loc    *time.Location
	layout string
	t      string
}{
	{20230613, 2023, 6, 13, time.UTC, DateTextUS, "June 13, 2023"},
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
		txt, err := tt.dt8.Format(tt.layout, tt.loc)
		if txt != tt.t {
			t.Errorf("timeutil.DateTime8.Format(\"%s\") val (%d) Error: err (%s)",
				tt.layout, tt.dt8, err.Error())
		}
		if txt != tt.t {
			t.Errorf("timeutil.DateTime8.Format(\"%s\") val (%d) Mismatch: want (%s) got (%s)",
				tt.layout, tt.dt8, tt.t, txt)
		}

		dt, err := tt.dt8.Time(tt.loc)
		if err != nil {
			t.Errorf("timeutil.DateTime8.Time() val (%d) Error: err (%s)",
				tt.dt8, err.Error())
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
		got, err := DT8ParseString(time.RFC3339, tt.v)
		if err != nil {
			t.Errorf("Dt8ForString(%v): want %v, error %v", tt.v, tt.want, err)
		}
		if got != tt.want {
			t.Errorf("Dt8ForString(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}

type testDateTime8UnmarshalStruct struct {
	DateTime8 DateTime8
}

var datetime8UnmarshalJSONTests = []struct {
	v     string
	isErr bool
	want  DateTime8
}{
	{`{"DateTime8":null}`, true, DateTime8(int32(0))},
	{`{"DateTime8":"string"}`, true, DateTime8(int32(0))},
	{`{"DateTime8":20230632}`, true, DateTime8(int32(0))},
	{`{"DateTime8":20230630}`, false, DateTime8(int32(20230630))},
}

func TestDateTime8UnmarshalJSON(t *testing.T) {
	for _, tt := range datetime8UnmarshalJSONTests {
		var w testDateTime8UnmarshalStruct
		err := json.Unmarshal([]byte(tt.v), &w)
		if err != nil {
			if !tt.isErr {
				t.Errorf("datetime8: json.Unmarshal(%s): error (%s)", tt.v, err.Error())
			}
		}
	}
}

var dtParseUint32sTests = []struct {
	yyyy uint32
	mm   uint32
	dd   uint32
	want DateTime8
}{
	{1582, 10, 15, 15821015},
	{1776, 7, 4, 17760704},
	{1896, 4, 6, 18960406},
	{2006, 1, 2, 20060102},
}

func TestDT8ParseUint32s(t *testing.T) {
	for _, tt := range dtParseUint32sTests {
		try, err := DT8ParseUint32s(tt.yyyy, tt.mm, tt.dd)
		if err != nil {
			t.Errorf("timeutil.DT8ParseUints(%d,%d,%d): error (%s)", tt.yyyy, tt.mm, tt.dd, err.Error())
		}
		if try != tt.want {
			t.Errorf("timeutil.DT8ParseUints(%d,%d,%d): want (%d), got (%d)", tt.yyyy, tt.mm, tt.dd, tt.want, try)
		}
	}
}
