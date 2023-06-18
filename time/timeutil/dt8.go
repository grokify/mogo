package timeutil

import (
	"fmt"
	"strconv"
	"time"
)

// DateTime8 represents a datetime `int32` value in the `yyyymmdd` format.
type DateTime8 int32

func (dt8 DateTime8) Format(layout string, loc *time.Location) (string, error) {
	dt, err := dt8.Time(loc)
	if err != nil {
		return "", err
	}
	return dt.Format(layout), nil
}

func (dt8 DateTime8) Split() (int32, int32, int32) {
	year := dt8 / 10000
	month := int(dt8/100) - (int(year) * 100)
	day := int(dt8) - (int(year) * 10000) - (month * 100)
	return int32(year), int32(month), int32(day)
}

/*
// DurationForNowSubDt8 returns a duartion struct between a Dt8 value and the current time.
func (dt8 DateTime8) DurationForNowSubDT8(dt8 int32) (time.Duration, error) {
	t, err := TimeForDT8(dt8)
	if err != nil {
		var d time.Duration
		return d, err
	}
	now := time.Now()
	return now.Sub(t), nil
}
*/

// Sub returns the duration dt8-u. If the result exceeds the maximum (or minimum) value that can be stored in a Duration, the maximum (or minimum) duration will be returned. To compute dt8-d for a duration d, use t.Add(-d).
func (dt8 DateTime8) SubTime(u time.Time, loc *time.Location) (time.Duration, error) {
	t, err := dt8.Time(loc)
	if err != nil {
		return 0, err
	}
	return t.Sub(u), nil
}

/*
func (dt8 DateTime8) Parse() time.Time {
	y, m, d := dt8.Split()
	return time.Date(int(y), time.Month(int(m)), int(d), 0, 0, 0, 0, time.UTC)
}
*/

func (dt8 DateTime8) Time(loc *time.Location) (time.Time, error) {
	dt, err := time.Parse(DT8, strconv.FormatInt(int64(dt8), 10))
	if loc == nil || (loc == dt.Location()) {
		return dt, err
	}
	return time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, loc), nil
}

func (dt8 DateTime8) Validate() error {
	_, err := dt8.Time(time.UTC)
	if err != nil {
		return err
	}
	return nil
}

// DT8ParseString returns a Dt8 value given a layout and value to parse to time.Parse.
func DT8ParseString(layout, value string) (DateTime8, error) {
	dt8 := DateTime8(int32(0))
	t, err := time.Parse(layout, value)
	if err != nil {
		return dt8, err
	}
	return NewTimeMore(t, 0).DT8(), nil
}

// DT8ParseUnts returns a Dt8 value for year, month, and day.
func DT8ParseUnts(yyyy, mm, dd uint) (DateTime8, error) {
	dt8String := fmt.Sprintf("%04d%02d%02d", yyyy, mm, dd)
	dt8Int, err := strconv.ParseInt(dt8String, 10, 32)
	if err != nil {
		panic(err)
	}
	dt8 := DateTime8(int32(dt8Int))
	return dt8, dt8.Validate()
}
