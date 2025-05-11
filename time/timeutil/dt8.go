package timeutil

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/mogo/strconv/strconvutil"
)

var ErrDateTime8OutOfBounds = errors.New("datetime8: time.Time is out of bounds")

// DateTime8 represents a datetime `int32` value in the `yyyymmdd` format. It supports
// dates from 1000-01-01 to 9999-12-31.
type DateTime8 uint32

func (dt8 DateTime8) Format(layout string, loc *time.Location) (string, error) {
	dt, err := dt8.Time(loc)
	if err != nil {
		return "", err
	}
	return dt.Format(layout), nil
}

func (dt8 DateTime8) Split() (year uint32, month uint32, day uint32) {
	year = uint32(dt8) / 10000
	month = uint32(dt8)/100 - (year * 100)
	day = uint32(dt8) - (year * 10000) - (month * 100)
	return
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

func (dt8 DateTime8) Time(loc *time.Location) (time.Time, error) {
	dt, err := time.Parse(DT8, strconv.Itoa(int(dt8)))
	if err != nil || loc == nil || (loc == dt.Location()) {
		return dt, err
	}
	return time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, loc), nil
}

var ErrDateTime8Invalid = errors.New("timeutil.datetime8: invalid value")

type DateTime8UnmarshalError struct {
	Msg string
}

func (e *DateTime8UnmarshalError) Error() string {
	return fmt.Sprintf("timeutil.datetime8: unmarshal error (%s)", e.Msg)
}

func (dt8 *DateTime8) UnmarshalJSON(data []byte) error {
	s := strings.TrimSpace(string(data))
	if len(s) == 0 || s == "0" {
		*dt8 = 0
		return nil
	}
	i, err := strconvutil.Atou32(s)
	if err != nil {
		return &DateTime8UnmarshalError{Msg: err.Error()}
	}
	d8 := DateTime8(i)
	err = d8.Validate()
	if err != nil {
		return err
	}
	*dt8 = d8
	return nil
}

func (dt8 DateTime8) Validate() error {
	_, err := dt8.Time(time.UTC)
	if err != nil {
		return err
	}
	return nil
}

// DT8ParseString returns a `DateTime8` value given a layout and value to parse to time.Parse.
func DT8ParseString(layout, value string) (DateTime8, error) {
	dt8 := DateTime8(int32(0))
	t, err := time.Parse(layout, value)
	if err != nil {
		return dt8, err
	}
	return NewTimeMore(t, 0).DT8()
}

// DT8ParseUints returns a `DateTime8` value for year, month, and day.
func DT8ParseUint32s(yyyy, mm, dd uint32) (DateTime8, error) {
	dt8 := DateTime8(yyyy*100*100 + mm*100 + dd)
	return dt8, dt8.Validate()
}

func dt8TimeInbounds(t time.Time) bool {
	if t.Year() < 1000 || t.Year() > 9999 {
		return false
	}
	return true
}
