package timeutil

import (
	"testing"
	"time"
)

var quarterStartAddTests = []struct {
	input            string
	quarterStart     string
	quarterDiffCount int
	quarterDiffTime  string
}{
	{"2017-01-01T00:00:00Z", "2017-01-01T00:00:00Z", 1, "2017-04-01T00:00:00Z"},
	{"2017-01-01T00:00:00Z", "2017-01-01T00:00:00Z", -1, "2016-10-01T00:00:00Z"},
	{"2017-01-01T00:00:00Z", "2017-01-01T00:00:00Z", -2, "2016-07-01T00:00:00Z"},
	{"2017-06-30T23:00:00Z", "2017-04-01T00:00:00Z", 2, "2017-10-01T00:00:00Z"},
	{"2017-07-06T00:00:00Z", "2017-07-01T00:00:00Z", 3, "2018-04-01T00:00:00Z"}}

func TestQuarterStartAdd(t *testing.T) {
	for _, tt := range quarterStartAddTests {
		dt, err := time.Parse(time.RFC3339, tt.input)
		if err != nil {
			t.Errorf("time.Parse(%v): want %v, err %v", tt.input, tt.quarterStart, err)
			continue
		}
		got := QuarterStart(dt)
		if got.Format(time.RFC3339) != tt.quarterStart {
			t.Errorf("QuarterStart(%v): want %v, got %v", tt.input, tt.quarterStart, got.Format(time.RFC3339))
		}
		diff := QuarterAdd(dt, tt.quarterDiffCount)
		diffRFC := diff.Format(time.RFC3339)
		if diffRFC != tt.quarterDiffTime {
			t.Errorf("QuarterAdd(%s, %d): want [%s], got [%v]", tt.input, tt.quarterDiffCount, tt.quarterDiffTime, diffRFC)
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
