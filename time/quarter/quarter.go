package quarter

import (
	"time"

	"github.com/grokify/mogo/math/mathutil"
	"github.com/grokify/mogo/type/number"
)

// YearQuarterToQuarterContinuous converts a year and quarter to
// a continuous quarter integer. This is useful when an even
// even spacing between months is desired, such as with
// charting x-axis values.
func YearQuarterToQuarterContinuous(year, quarter uint32) uint32 {
	return year*4 + quarter
}

// QuarterContinuousToYearQuarter converts a continuous quarter value (e.g. numerof months from year 0).
func QuarterContinuousToYearQuarter(qc uint32) (uint32, uint32) {
	quotient, remainder := mathutil.Divide(qc-1, 4)
	return quotient, remainder + 1
}

// TimeToQuarterContinuous converts a `time.Time` value to a continuous quarter.
func TimeToQuarterContinuous(t time.Time) (uint32, error) {
	if y, err := number.Itou32(t.Year()); err != nil {
		return 0, err
	} else if m, err := number.Itou32(int(t.Month())); err != nil {
		return 0, err
	} else {
		return YearQuarterToQuarterContinuous(y, MonthToQuarter(m)), nil
	}
}

func MonthToQuarter(month uint32) uint32 {
	if month <= 3 {
		return uint32(1)
	} else if month <= 6 {
		return uint32(2)
	} else if month <= 9 {
		return uint32(3)
	}
	return uint32(4)
}

// QuarterContinuousToTime converts a continuous quarter value to a `time.Time` value.
func QuarterContinuousToTime(qc uint32) time.Time {
	year, quarter := QuarterContinuousToYearQuarter(qc)
	month := 1
	switch quarter {
	case 2:
		month = 4
	case 3:
		month = 7
	case 4:
		month = 10
	}
	return time.Date(
		int(year), time.Month(month), 1,
		0, 0, 0, 0, time.UTC)
}

func QuarterContinuousIsQuarterStart(qc uint32) bool {
	month := QuarterContinuousToTime(qc).Month()
	return month == 1 || month == 4 || month == 7 || month == 10
}

func QuarterContinuousIsYearStart(qc uint32) bool {
	return QuarterContinuousToTime(qc).Month() == 1
}
