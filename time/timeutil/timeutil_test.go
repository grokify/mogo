package timeutil

import (
	"testing"
)

var dt6ForDT14Tests = []struct {
	v    int
	want int
}{
	{20060102150405, 20060102}}

func TestDT6ForDT14(t *testing.T) {
	for _, tt := range dt6ForDT14Tests {
		got := DT6ForDT14(tt.v)
		if got != tt.want {
			t.Errorf("Dt6ForDt14(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}

var addMonthsTests = []struct {
	year      int
	month     int
	add       int
	wantYear  int
	wantMonth int
}{
	{2025, 1, -1, 2024, 12},
	{2025, 1, 1, 2025, 2},
	{2025, 12, 1, 2026, 1},
}

func TestAddMonths(t *testing.T) {
	for _, tt := range addMonthsTests {
		gotYear, gotMonth := AddMonths(tt.year, tt.month, tt.add)
		if gotYear != tt.wantYear || gotMonth != tt.wantMonth {
			t.Errorf("AddMonths(%d,%d,%d): want (%d,%d), got (%d,%d)",
				tt.year, tt.month, tt.add,
				tt.wantYear, tt.wantMonth,
				gotYear, gotMonth)
		}
	}
}

/*
var fromToTests = []struct {
	v    string
	want string
}{
	{"Wed, 25 May 2016 11:07:20 +0000", "2016-05-25T11:07:20Z"}}

func FromToTest(t *testing.T) {
	for _, tt := range fromToTests {
		got, err := FromTo(tt.v, time.RFC1123Z, time.RFC3339)
		if err != nil {
			t.Errorf("FromTo(%v): want %v, error ", tt.v, tt.want, err)
		}
		if got != tt.want {
			t.Errorf("FromTo(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}
*/
