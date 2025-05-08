package month

import (
	"fmt"
	"time"

	"github.com/grokify/base36"
	"github.com/grokify/mogo/math/mathutil"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/number"
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

func YearMonthBase36(yyyy, mm uint32) string {
	return fmt.Sprintf("%04s", base36.Encode(uint64(yyyy*100+mm)))
}

func YearMonthBase36Time(t time.Time) (string, error) {
	if y, err := number.Itou32(t.Year()); err != nil {
		return "", err
	} else if m, err := number.Itou32(int(t.Month())); err != nil {
		return "", err
	} else {
		return YearMonthBase36(y, m), nil
	}
}

// MonthStart allows you to add/subtract months resulting
// in the first day of each month while avoiding Go's
// `AddDate` normalization where "adding one month to
// October 31 yields December 1, the normalized form for
// November 31." Setting `deltaMonths` to 0 indicates the
// current month.
func MonthStart(t time.Time, deltaMonths int) time.Time {
	year := t.Year()
	month := int(t.Month())
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
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, t.Location())
}

// YearMonthToMonthContinuous converts a year and month to
// a continuous month integer. This is useful when an even
// even spacing between months is desired, such as with
// charting x-axis values.
func YearMonthToMonthContinuous(year, month uint32) uint32 {
	return year*12 + month
}

// MonthContinuousToYearMonth converts a continuous month value (e.g. number of months from year 0).
func MonthContinuousToYearMonth(mc uint32) (uint32, uint32) {
	quotient, remainder := mathutil.Divide(mc-1, 12)
	return quotient, remainder + 1
}

// TimeToMonthContinuous converts a `time.Time` value to a continuous month.
func TimeToMonthContinuous(t time.Time) (uint32, error) {
	t = t.UTC()
	if y, err := number.Itou32(t.Year()); err != nil {
		return 0, err
	} else if m, err := number.Itou32(int(t.Month())); err != nil {
		return 0, err
	} else {
		return YearMonthToMonthContinuous(y, m), nil
	}
}

// MonthContinuousToTime converts a continuous month value to a `time.Time` value.
func MonthContinuousToTime(mc uint32) time.Time {
	year, month := MonthContinuousToYearMonth(mc)
	return time.Date(
		int(year), time.Month(int(month)), 1,
		0, 0, 0, 0, time.UTC)
}

func MonthContinuousIsQuarterStart(mc uint32) bool {
	month := MonthContinuousToTime(mc).Month()
	return month == 1 || month == 4 || month == 7 || month == 10
}

func MonthContinuousIsYearStart(mc uint32) bool {
	return MonthContinuousToTime(mc).Month() == 1
}

// TimesMonthStarts returns time series of months given start and end input times.
func TimesMonthStarts(times ...time.Time) timeutil.Times {
	timesStarts := timeutil.Times{}
	if len(times) == 0 {
		return timesStarts
	}
	min, max := timeutil.SliceMinMax(times)
	minMonth := timeutil.NewTimeMore(min, 0).MonthStart()
	maxMonth := timeutil.NewTimeMore(max, 0).MonthStart()
	curMonth := minMonth
	for curMonth.Before(maxMonth) || curMonth.Equal(maxMonth) {
		timesStarts = append(timesStarts, curMonth)
		// curMonth = timeutil.TimeDt6AddNMonths(curMonth, 1)
		curMonth = curMonth.AddDate(0, 1, 0)
	}
	return timesStarts
}
