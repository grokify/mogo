package timeutil

import (
	"fmt"
	"time"

	"github.com/grokify/base36"
)

func DayofmonthToEnglish(i uint16) string {
	days := []string{
		"zeorth",
		"first", "second", "third", "fourth", "fifth",
		"sixth", "seventh", "eighth", "ninth", "tenth",
		"eleventh", "twelfth", "thirteenth", "fourteenth", "fifteenth",
		"sixteenth", "seventeenth", "eighteenth", "nineteenth", "twentieth",
	}
	tenZero := []string{"tenth", "twentieth", "thirtieth", "fourtieth", "fiftieth"}
	tenPlus := []string{"ten", "twenty", "thirty", "forty", "fifty"}
	if i < 21 {
		return days[i]
	}
	if i > 59 {
		panic("E_OUT_OF_RANGE")
	}
	quotient, remainder := i/10, i%10
	if remainder == 0 {
		return tenZero[quotient-1]
	}
	return tenPlus[quotient-1] + " " + days[remainder]
}

func YearMonthBase36(yyyy, mm uint64) string {
	return fmt.Sprintf("%04s", base36.Encode(yyyy*100+mm))
}

func YearMonthBase36Time(dt time.Time) string {
	return YearMonthBase36(uint64(dt.Year()), uint64(dt.Month()))
}

// MonthBegin allows you to add/subtract months resulting
// in the first day of each month while avoiding Go's
// `AddDate` normalization where "adding one month to
// October 31 yields December 1, the normalized form for November 31."
func MonthBegin(dt time.Time, deltaMonths int) time.Time {
	dt = dt.UTC()
	year := dt.Year()
	month := int(dt.Month())
	if deltaMonths > 0 {
		for i := 0; i < deltaMonths; i++ {
			if month == 12 {
				month = 1
				year++
			} else {
				month++
			}
		}
	} else if deltaMonths < 0 {
		for i := 0; i > deltaMonths; i-- {
			if month == 1 {
				month = 12
				year--
			} else {
				month--
			}
		}
	}
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
}
