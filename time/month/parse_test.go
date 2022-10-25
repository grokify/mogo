package month

import (
	"strings"
	"testing"
	"time"
)

var parseMonthTests = []struct {
	haystack    []string
	value       string
	insensiitve bool
	month       time.Month
	err         error
}{
	{[]string{MonthsEnAbbr3}, "Jul", false, time.July, nil},
	{[]string{MonthsEnFull}, "   september   ", true, time.September, nil},
	{[]string{MonthsEnAbbr3}, "oct", true, time.October, nil},
	{[]string{MonthsEnAbbr3}, "Oct", false, time.October, nil},
}

func TestParseMonth(t *testing.T) {
	for _, tt := range parseMonthTests {
		monthTry, err := ParseMonth(tt.haystack, tt.value, tt.insensiitve)
		if err != nil {
			t.Errorf(`time.month.Parse("%s". "%s", %v) got error [%v]`,
				strings.Join(tt.haystack, ","), tt.value, tt.insensiitve, err.Error())
		}
		if monthTry != tt.month {
			t.Errorf(`time.month.Parse("%s", "%s", %v) Expected [%s] Got [%s]`,
				strings.Join(tt.haystack, ","), tt.value, tt.insensiitve, tt.month.String(), monthTry.String())
		}
	}
}
