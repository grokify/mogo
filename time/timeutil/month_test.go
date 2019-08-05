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
