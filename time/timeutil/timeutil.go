// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"fmt"
	"strconv"
	"time"
)

// TimeOpts represnts a struct for `time.Date`.
type TimeOpts struct {
	Year       int
	Month      int
	Day        int
	Hour       int
	Minute     int
	Second     int
	Nanosecond int
	Location   *time.Location
}

// Time returns a `time.Time` struct. If no `Location`
// is set, `time.UTC` is used.
func (opts TimeOpts) Time() time.Time {
	if opts.Location == nil {
		opts.Location = time.UTC
	}
	return time.Date(
		opts.Year,
		time.Month(opts.Month),
		opts.Day,
		opts.Hour,
		opts.Minute,
		opts.Second,
		opts.Nanosecond,
		opts.Location)
}

func MustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

// UnixToDay converts an epoch in seconds to a time.Time for the day.
func UnixToDay(epoch int64) time.Time {
	return NewTimeMore(time.Unix(epoch, 0).UTC(), 0).DayStart()
}

// DT6 returns the Dt6 value for time.Time.
func (tm TimeMore) DT6() int32 {
	return int32(tm.time.Year()*100 + int(tm.time.Month()))
}

// DT6ForDT14 returns the Dt6 value for Dt14.
func DT6ForDT14(dt14 int64) int32 {
	return int32(float64(dt14) / float64(1000000))
}

// TimeForDT6 returns a time.Time value given a Dt6 value.
func TimeForDT6(dt6 int32) (time.Time, error) {
	return time.Parse(DT6, strconv.FormatInt(int64(dt6), 10))
}

func DT6Parse(dt6 int32) (int16, int8) {
	year := dt6 / 100
	month := int(dt6) - (int(year) * 100)
	return int16(year), int8(month)
}

func DT6Prev(dt6 int32) int32 {
	year, month := DT6Parse(dt6)
	if month == 1 {
		month = 12
		year = year - 1
	} else {
		month = month - 1
	}
	return int32(year)*100 + int32(month)
}

func DT6Next(dt6 int32) int32 {
	year, month := DT6Parse(dt6)
	if month == 12 {
		month = 1
		year++
	} else {
		month++
	}
	return int32(year)*100 + int32(month)
}

func TimeDT6AddNMonths(t time.Time, numMonths int) time.Time {
	if numMonths == 0 {
		return t
	} else if numMonths < 0 {
		return timeDT6SubNMonths(t, uint(-numMonths))
	}
	dt6 := NewTimeMore(t, 0).DT6()
	for i := 0; i < numMonths; i++ {
		dt6 = DT6Next(dt6)
	}
	dt6NextMonth, err := TimeForDT6(dt6)
	if err != nil {
		panic(fmt.Sprintf("Cannot find next month for time: %v\n", t.Format(time.RFC3339)))
	}
	return dt6NextMonth
}

func timeDT6SubNMonths(t time.Time, numMonths uint) time.Time {
	if numMonths == 0 {
		return t
	}
	dt6 := NewTimeMore(t, 0).DT6()
	for i := uint(0); i < numMonths; i++ {
		dt6 = DT6Prev(dt6)
	}
	dt6NextMonth, err := TimeForDT6(dt6)
	if err != nil {
		panic(fmt.Sprintf("Cannot find next month for time: %v\n", t.Format(time.RFC3339)))
	}
	return dt6NextMonth
}

func TimeDT4AddNYears(t time.Time, numYears int) time.Time {
	return time.Date(t.Year()+numYears, time.January, 1, 0, 0, 0, 0, t.Location())
}

func DT6MinMaxSlice(minDt6 int32, maxDt6 int32) []int32 {
	if maxDt6 < minDt6 {
		tmpDt6 := maxDt6
		maxDt6 = minDt6
		minDt6 = tmpDt6
	}
	dt6Range := []int32{}
	curDt6 := minDt6
	for curDt6 < maxDt6+1 {
		dt6Range = append(dt6Range, curDt6)
		curDt6 = DT6Next(curDt6)
	}
	return dt6Range
}

/*
// Dt8Now returns Dt8 value for the current time.
func Dt8Now() int32 {
	return Dt8ForTime(time.Now())
}

// DT8ForString returns a Dt8 value given a layout and value to parse to time.Parse.
func DT8ForString(layout, value string) (int32, error) {
	dt8 := int32(0)
	t, err := time.Parse(layout, value)
	if err == nil {
		dt8 = NewTimeMore(t, 0).DT8()
	}
	return dt8, err
}

// DT8ForInts returns a Dt8 value for year, month, and day.
func DT8ForInts(yyyy, mm, dd int) int32 {
	sDt8 := fmt.Sprintf("%04d%02d%02d", yyyy, mm, dd)
	iDt8, err := strconv.ParseInt(sDt8, 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(iDt8)
}

// TimeForDT8 returns a time.Time value given a Dt8 value.
func TimeForDT8(dt8 int32) (time.Time, error) {
	return time.Parse(DT8, strconv.FormatInt(int64(dt8), 10))
}

// Dt14Now returns a Dt14 value for the current time.
func Dt14Now() int64 {
	return Dt14ForTime(time.Now())
}
*/

// Dt8ForTime returns a `DateTime8` value given a time struct.
func (tm TimeMore) DT8() (DateTime8, error) {
	if !dt8TimeInbounds(tm.time) {
		return 0, ErrDateTime8OutOfBounds
	}
	s := tm.time.Format(DT8)
	iDt8, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return DateTime8(iDt8), nil
}

// DT14ForString returns a DT14 value given a layout and value to parse to time.Parse.
func DT14ForString(layout, value string) (int64, error) {
	dt14 := int64(0)
	t, err := time.Parse(layout, value)
	if err == nil {
		dt14 = NewTimeMore(t, 0).DT14()
	}
	return dt14, err
}

// DT14ForInts returns a Dt8 value for a UTC year, month, day, hour, minute and second.
func DT14ForInts(yyyy, mm, dd, hr, mn, dy int) int64 {
	sDt14 := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", yyyy, mm, dd, hr, mn, dy)
	iDt14, err := strconv.ParseInt(sDt14, 10, 64)
	if err != nil {
		panic(err)
	}
	return iDt14
}

// Dt14ForTime returns a Dt14 value given a time.Time struct.
func (tm TimeMore) DT14() int64 {
	s := tm.time.Format(DT14)
	iDt14, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return iDt14
}

// TimeForDT14 returns a time.Time value given a Dt14 value.
func TimeForDT14(dt14 int64) (time.Time, error) {
	return time.Parse(DT14, strconv.FormatInt(dt14, 10))
}

/*
// WeekStart takes a time.Time object and a week start day
// in the time.Weekday format.
func WeekStart(dt time.Time, dow time.Weekday) (time.Time, error) {
	return TimeDeltaDowInt(dt.UTC(), int(dow), -1, true, true)
}

// MonthStart returns a time.Time for the beginning of the
// month in UTC time.
func MonthStart(dt time.Time) time.Time {
	dt = dt.UTC()
	return time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// QuarterEnd returns a time.Time for the end of the
// quarter by second in UTC time.
func QuarterEnd(dt time.Time) time.Time {
	qs := QuarterStart(dt.UTC())
	qn := TimeDt6AddNMonths(qs, 3)
	return time.Date(qn.Year(), qn.Month(), 0, 23, 59, 59, 0, time.UTC)
}

// YearStart returns a a time.Time for the beginning of the year
// in UTC time.
func YearStart(dt time.Time) time.Time {
	return time.Date(dt.UTC().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
}

// YearEnd returns a a time.Time for the end of the year in UTC time.
func YearEnd(dt time.Time) time.Time {
	return time.Date(dt.UTC().Year(), time.December, 31, 23, 59, 59, 999999999, time.UTC)
}

func NextYearStart(dt time.Time) time.Time {
	return time.Date(dt.UTC().Year()+1, time.January, 1, 0, 0, 0, 0, time.UTC)
}

func IsYearStart(t time.Time) bool {
	t = t.UTC()
	if t.Nanosecond() == 0 &&
		t.Second() == 0 &&
		t.Minute() == 0 &&
		t.Hour() == 0 &&
		t.Day() == 1 &&
		t.Month() == time.January {
		return true
	}
	return false
}

func ToMonthStart(t time.Time) time.Time {
	return time.Date(
		t.Year(), t.Month(), 1,
		0, 0, 0, 0,
		t.Location())
}
*/
