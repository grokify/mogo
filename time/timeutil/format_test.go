package timeutil

import (
	"encoding/json"
	"testing"
	"time"
)

var dmyhm2ParseTests = []struct {
	v    string
	want string
}{
	{"02:01:06 15:04", "2006-01-02T15:04:00Z"},
}

// TestDMYHM2ParseTime ensures timeutil.DateDMYHM2 is parsed to GMT timezone.
func TestDMYHM2ParseTime(t *testing.T) {
	for _, tt := range dmyhm2ParseTests {
		got, err := FromTo(tt.v, DateDMYHM2, time.RFC3339)
		if err != nil {
			t.Errorf("time.Parse(DateDMYHM2) Error: with %v, want %v, err %v", tt.v, tt.want, err)
		}
		if got != tt.want {
			t.Errorf("time.Parse(\"%v\", DateDMYHM2) Mismatch: want %v, got %v", tt.v, tt.want, got)
		}
	}
}

var rfc3339YMDTimeTests = []struct {
	v    string
	want string
}{
	{`{"MyTime":"2001-02-03"}`, `{"MyTime":"2001-02-03"}`},
	{`{"MyTime":"0001-01-01"}`, `{"MyTime":"0001-01-01"}`},
	{`{"MyTime":""}`, `{"MyTime":"0001-01-01"}`},
	{`{}`, `{"MyTime":"0001-01-01"}`}}

type myStruct struct {
	MyTime RFC3339YMDTime
}

func TestRfc3339YMDTime(t *testing.T) {
	for _, tt := range rfc3339YMDTimeTests {
		my := myStruct{}
		//fmt.Println(tt.v)
		err := json.Unmarshal([]byte(tt.v), &my)
		if err != nil {
			t.Errorf("Rfc3339YMDTime.Unmarshal: with %v, want %v, err %v", tt.v, tt.want, err)
		}

		bytes, err := json.Marshal(my)
		if err != nil {
			t.Errorf("Rfc3339YMDTime.Marshal: with %v, want %v, err %v", tt.v, tt.want, err)
		}

		got := string(bytes)

		if got != tt.want {
			t.Errorf("Rfc3339YMDTime(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}

var offsetFormatTests = []struct {
	input    int
	useColon bool
	useZ     bool
	want     string
}{
	{0, false, false, "+0000"},
	{0, true, false, "+00:00"},
	{0, true, true, "Z"},
	{400, false, false, "+0400"},
	{-400, false, false, "-0400"},
	{530, false, false, "+0530"},
	{-530, false, false, "-0530"},
	{700, true, false, "+07:00"},
	{-700, true, false, "-07:00"},
	{845, true, false, "+08:45"},
	{-845, true, false, "-08:45"},
}

func TestOffsetFormat(t *testing.T) {
	for _, tt := range offsetFormatTests {
		got := OffsetFormat(tt.input, tt.useColon, tt.useZ)
		if got != tt.want {
			t.Errorf("time.OffsetFormat(\"%v\",%v,%v) Mismatch: want [%v], got [%v]", tt.input, tt.useColon, tt.useZ, tt.want, got)
		}
	}
}

var isDTXTests = []struct {
	v    int
	want string
}{
	{2004, LayoutNameDT4},
	{200401, LayoutNameDT6},
}

func TestIsDTX(t *testing.T) {
	for _, tt := range isDTXTests {
		got, err := IsDTX(tt.v)
		if err != nil {
			t.Errorf("time.IsDTX(%d) Error: want (%s), err (%v)", tt.v, tt.want, err)
		}
		if got != tt.want {
			t.Errorf("time.IsDTX(%d) Mismatch: want (%s), got (%s)", tt.v, tt.want, got)
		}
	}
}

var formatsTests = []struct {
	v      string
	format string
	want   string
}{
	{"2023-06-18T00:00:00Z", DateTextUS, "June 18, 2023"},
	{"2023-06-18T00:00:00Z", DateTextUSAbbr3, "Jun 18, 2023"},
}

func TestFormats(t *testing.T) {
	for _, tt := range formatsTests {
		dt, err := time.Parse(time.RFC3339, tt.v)
		if err != nil {
			t.Errorf("time.Parse(time.RFC3339, %s) Error: err (%s)", tt.v, err.Error())
		}
		got := dt.Format(tt.format)
		if got != tt.want {
			t.Errorf("time.Format(%s) dt (%s) Mismatch: want (%s), got (%s)", tt.format, tt.v, tt.want, got)
		}
	}
}
