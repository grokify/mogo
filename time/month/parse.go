package month

import (
	"errors"
	"strings"
	"time"

	"github.com/grokify/mogo/type/stringsutil"
)

const (
	MonthsEnAbbr3 = "Jan,Feb,Mar,Apr,May,Jun,Jul,Aug,Sep,Oct,Nov,Dec"
	MonthsEnFull  = "January,February,March,April,May,June,July,August,September,October,November,December"
)

var (
	ErrMonthsFormatInvalid = errors.New("invalid number of month elements")
	ErrMonthNotFound       = errors.New("month not found")
)

func Parse(format []string, value string, insensitive bool) (time.Month, error) {
	months := format
	if len(months) == 1 {
		months = strings.Split(format[0], ",")
	}
	months = stringsutil.SliceCondenseSpace(months, true, false)
	if len(months) != 12 {
		return time.January, ErrMonthsFormatInvalid
	}
	value = strings.TrimSpace(value)
	if insensitive {
		value = strings.ToLower(value)
	}

	for i, monthTry := range months {
		if (insensitive && strings.ToLower(monthTry) == value) ||
			monthTry == value {
			return time.Month(i + 1), nil
		} else if monthTry == value {
			return time.Month(i + 1), nil
		}
	}
	return time.January, ErrMonthNotFound
}
