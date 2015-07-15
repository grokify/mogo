package timeutil

import (
	"fmt"
	"strconv"
	"time"
)

const (
	DT14 = "20060102150405"
	DT8  = "20060102"
)

func IsGreaterThan(timeLeft time.Time, timeRight time.Time) bool {
	durDelta := timeLeft.Sub(timeRight)
	if durZero, _ := time.ParseDuration("0ns"); durDelta > durZero {
		return true
	}
	return false
}

func IsLessThan(timeLeft time.Time, timeRight time.Time) bool {
	durDelta := timeLeft.Sub(timeRight)
	if durZero, _ := time.ParseDuration("0ns"); durDelta < durZero {
		return true
	}
	return false
}

func Dt8Now() int32 {
	return Dt8ForTime(time.Now())
}

func Dt8ForString(layout, value string) (int32, error) {
	dt8 := int32(0)
	t, err := time.Parse(layout, value)
	if err == nil {
		dt8 = Dt8ForTime(t)
	}
	return dt8, err
}

func Dt8ForInts(yyyy int, mm int, dd int) int32 {
	sDt8 := fmt.Sprintf("%04d%02d%02d", yyyy, mm, dd)
	iDt8, _ := strconv.ParseInt(sDt8, 10, 32)
	return int32(iDt8)
}

func Dt8ForTime(t time.Time) int32 {
	u := t.UTC()
	s := u.Format(DT8)
	iDt8, _ := strconv.ParseInt(s, 10, 32)
	return int32(iDt8)
}

func TimeForDt8(dt8 int32) (time.Time, error) {
	return time.Parse(DT8, strconv.FormatInt(int64(dt8), 10))
}

func DurationForNowSubDt8(dt8 int32) (time.Duration, error) {
	t, err := TimeForDt8(dt8)
	if err != nil {
		var d time.Duration
		return d, err
	}
	now := time.Now()
	return now.Sub(t), nil
}

func Dt14Now() int64 {
	return Dt14ForTime(time.Now())
}

func Dt14ForString(layout, value string) (int64, error) {
	dt14 := int64(0)
	t, err := time.Parse(layout, value)
	if err == nil {
		dt14 = Dt14ForTime(t)
	}
	return dt14, err
}

func Dt14ForInts(yyyy int, mm int, dd int, hr int, mn int, dy int) int64 {
	sDt14 := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", yyyy, mm, dd, hr, mn, dy)
	iDt14, _ := strconv.ParseInt(sDt14, 10, 64)
	return int64(iDt14)
}

func Dt14ForTime(t time.Time) int64 {
	u := t.UTC()
	s := u.Format(DT14)
	iDt14, _ := strconv.ParseInt(s, 10, 64)
	return int64(iDt14)
}

func TimeForDt14(dt14 int64) (time.Time, error) {
	return time.Parse(DT14, strconv.FormatInt(dt14, 10))
}
