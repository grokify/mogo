package timeutil

import (
	"testing"
	"time"
)

var quarterStartTests = []struct {
	v    string
	want string
}{
	{"2017-01-01T00:00:00Z", "2017-01-01T00:00:00Z"},
	{"2017-06-30T23:00:00Z", "2017-04-01T00:00:00Z"},
	{"2017-07-06T00:00:00Z", "2017-07-01T00:00:00Z"}}

func TestQuarterStart(t *testing.T) {
	for _, tt := range quarterStartTests {
		dt, err := time.Parse(time.RFC3339, tt.v)
		if err != nil {
			t.Errorf("time.Parse(%v): want %v, err %v", tt.v, tt.want, err)
			continue
		}
		got := QuarterStart(dt)
		if got.Format(time.RFC3339) != tt.want {
			t.Errorf("QuarterStart(%v): want %v, got %v", tt.v, tt.want, got.Format(time.RFC3339))
		}
	}
}

var monthToQuarterTests = []struct {
	v    uint8
	want uint8
}{
	{1, 1}, {2, 1}, {3, 1},
	{4, 2}, {5, 2}, {6, 2},
	{7, 3}, {8, 3}, {9, 3},
	{10, 4}, {11, 4}, {12, 4}}

func TestMonthToQuarter(t *testing.T) {
	for _, tt := range monthToQuarterTests {
		got := MonthToQuarter(tt.v)
		if got != tt.want {
			t.Errorf("MonthToQuarter(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}

var quarterToMonthTests = []struct {
	v    uint8
	want uint8
}{
	{1, 1}, {2, 4}, {3, 7}, {4, 10}}

func TestQuarterToMonth(t *testing.T) {
	for _, tt := range quarterToMonthTests {
		got := QuarterToMonth(tt.v)
		if got != tt.want {
			t.Errorf("QuarterToMonth(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}
