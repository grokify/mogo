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

func Parse(months []string, value string, insensitive bool) (time.Month, error) {
	enums := months
	if len(enums) == 1 {
		enums = strings.Split(enums[0], ",")
	}
	enums = stringsutil.SliceCondenseSpace(enums, true, false)
	if len(enums) != 12 {
		return time.January, ErrMonthsFormatInvalid
	}
	value = strings.TrimSpace(value)
	if insensitive {
		value = strings.ToLower(value)
	}

	for i, monthTry := range enums {
		if (insensitive && strings.ToLower(monthTry) == value) ||
			monthTry == value {
			return time.Month(i + 1), nil
		}
	}
	return time.January, ErrMonthNotFound
}
