// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseDuration adds days (d), weeks (w), years (y).
func ParseDuration(s string) (time.Duration, error) {
	rx := regexp.MustCompile(`(?i)^\s*(-?\d+)(d|w|y)\s*$`)
	rs := rx.FindStringSubmatch(s)

	if len(rs) > 0 {
		zeroDuration, _ := time.ParseDuration("0s")
		quantity := rs[1]
		units := strings.ToLower(rs[2])
		i, err := strconv.Atoi(quantity)
		if err != nil {
			return zeroDuration, err
		}
		if units == "d" {
			s = fmt.Sprintf("%vs", i*DaySeconds)
		} else if units == "w" {
			s = fmt.Sprintf("%vs", i*WeekSeconds)
		} else if units == "y" {
			s = fmt.Sprintf("%vs", i*YearSeconds)
		} else {
			return zeroDuration, errors.New("timeutil.ParseDuration Parse Error")
		}
	}
	return time.ParseDuration(s)
}

func MustParseDuration(s string) time.Duration {
	dur, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return dur
}

func NowDeltaDuration(d time.Duration) time.Time {
	return time.Now().Add(d)
}

func NowDeltaParseDuration(s string) (time.Time, error) {
	d, err := ParseDuration(s)
	if err != nil {
		return time.Now(), err
	}
	return time.Now().Add(d), nil
}

// DurationForNowSubDt8 returns a duartion struct between a Dt8 value and the current time.
func DurationForNowSubDt8(dt8 int32) (time.Duration, error) {
	t, err := TimeForDt8(dt8)
	if err != nil {
		var d time.Duration
		return d, err
	}
	now := time.Now()
	return now.Sub(t), nil
}

func DurationStringMinutesSeconds(durationSeconds int64) (string, error) {
	if durationSeconds <= 0 {
		return "0 sec", nil
	}
	dur, err := time.ParseDuration(fmt.Sprintf("%vs", durationSeconds))
	if err != nil {
		return "", err
	}
	modSeconds := math.Mod(float64(durationSeconds), float64(60))
	if dur.Minutes() < 1 {
		return fmt.Sprintf("%v sec", modSeconds), nil
	}
	return fmt.Sprintf("%v min %v sec", int(dur.Minutes()), modSeconds), nil
}

// QuarterDuration returns a time.Duration representing the
// calendar quarter for the time provided.
func QuarterDuration(dt time.Time) time.Duration {
	start := QuarterStart(dt)
	end := NextQuarter(start)
	return end.Sub(start)
}

func SumDurations(durations ...time.Duration) time.Duration {
	seconds := int64(0)
	for _, dur := range durations {
		seconds += dur.Nanoseconds()
	}
	dur, _ := time.ParseDuration(fmt.Sprintf("%vns", seconds))
	return dur
}
