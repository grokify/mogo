package timeutil

import (
	"testing"
)

var dayofmonthToEnglishTests = []struct {
	v    uint16
	want string
}{
	{uint16(1), "first"},
	{uint16(2), "second"},
	{uint16(3), "third"},
	{uint16(19), "nineteenth"},
	{uint16(20), "twentieth"},
	{uint16(21), "twenty first"},
	{uint16(30), "thirtieth"},
	{uint16(31), "thirty first"},
}

func TestDayofmonthToEnglish(t *testing.T) {
	for _, tt := range dayofmonthToEnglishTests {
		got := DayofmonthToEnglish(tt.v)
		if got != tt.want {
			t.Errorf("DayofmonthToEnglish(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}

var yearMonthBase36Tests = []struct {
	year  uint64
	month uint64
	want  string
}{
	{uint64(1), uint64(1), "002T"},
	{uint64(500), uint64(3), "12KZ"},
	{uint64(1000), uint64(6), "255Y"},
	{uint64(2010), uint64(8), "4B3K"},
	{uint64(2019), uint64(8), "4BSK"},
	{uint64(2038), uint64(1), "4D95"},
	{uint64(9999), uint64(12), "LFJC"},
}

func TestYearMonthBase36(t *testing.T) {
	for _, tt := range yearMonthBase36Tests {
		got := YearMonthBase36(tt.year, tt.month)
		if got != tt.want {
			t.Errorf("YearMonthBase36(%v, %v): want %v, got %v", tt.year, tt.month, tt.want, got)
		}
	}
}
