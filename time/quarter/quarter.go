package quarter

import (
	"time"

	"github.com/grokify/mogo/math/mathutil"
)

// YearQuarterToQuarterContinuous converts a year and quarter to
// a continuous quarter integer. This is useful when an even
// even spacing between months is desired, such as with
// charting x-axis values.
func YearQuarterToQuarterContinuous(year, quarter uint64) uint64 {
	return year*4 + quarter
}

// QuarterContinuousToYearQuarter converts a continuous quarter value (e.g. numerof months from year 0).
func QuarterContinuousToYearQuarter(qc uint64) (uint64, uint64) {
	quotient, remainder := mathutil.DivideInt64(
		int64(qc-1), int64(4))
	return uint64(quotient), uint64(remainder + 1)
}

// TimeToQuarterContinuous converts a `time.Time` value to a continuous quarter.
func TimeToQuarterContinuous(t time.Time) uint64 {
	return YearQuarterToQuarterContinuous(
		uint64(t.Year()), MonthToQuarter(uint64(t.Month())))
}

func MonthToQuarter(month uint64) uint64 {
	if month <= 3 {
		return uint64(1)
	} else if month <= 6 {
		return uint64(2)
	} else if month <= 9 {
		return uint64(3)
	}
	return uint64(4)
}

// QuarterContinuousToTime converts a continuous quarter value to a `time.Time` value.
func QuarterContinuousToTime(qc uint64) time.Time {
	year, quarter := QuarterContinuousToYearQuarter(qc)
	month := 1
	if quarter == 2 {
		month = 4
	} else if quarter == 3 {
		month = 7
	} else if quarter == 4 {
		month = 10
	}
	return time.Date(
		int(year), time.Month(month), 1,
		0, 0, 0, 0, time.UTC)
}

func QuarterContinuousIsQuarterStart(qc uint64) bool {
	month := QuarterContinuousToTime(qc).Month()
	return month == 1 || month == 4 || month == 7 || month == 10
}

func QuarterContinuousIsYearStart(qc uint64) bool {
	return QuarterContinuousToTime(qc).Month() == 1
}
