package timeutil

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
)

// NewDurationSeconds returns a new `time.Duration` given a number of seconds.
func NewDurationSeconds(secs float64) time.Duration {
	nanos := int64(secs * float64(nanosPerSecond))
	dur, err := time.ParseDuration(strconv.Itoa(int(nanos)) + "ns")
	if err != nil {
		panic(err)
	}
	return dur
}

// NewDurationDays returns `time.Duration` given a number of days
func NewDurationDays(days uint16) time.Duration {
	durString := fmt.Sprintf("%dh", 24*days)
	dur, err := time.ParseDuration(durString)
	if err != nil {
		panic(err)
	}
	return dur
}

func NewDurationStrings(h, m, s string) (time.Duration, error) {
	return time.ParseDuration(fmt.Sprintf("%vh%vm%vs", h, m, s))
}

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
	end := QuarterAdd(start, 1)
	return end.Sub(start)
}

func SumDurations(durations ...time.Duration) time.Duration {
	ns := int64(0)
	for _, dur := range durations {
		ns += dur.Nanoseconds()
	}
	dur, _ := time.ParseDuration(fmt.Sprintf("%dns", ns))
	return dur
}

// SubDuration subtracts one duration from another and
// returns the result as a `time.Duration`.
func SubDuration(dur1, dur2 time.Duration) time.Duration {
	ns := dur1.Nanoseconds() - dur2.Nanoseconds()
	diff, err := time.ParseDuration(fmt.Sprintf("%dns", ns))
	if err != nil {
		panic("err")
	}
	return diff
}

func DurationIsZero(dur time.Duration) bool {
	return dur.Nanoseconds() == 0
}

func DurationZero() time.Duration {
	dur, _ := time.ParseDuration("0s")
	return dur
}

func MaxDuration(durs []time.Duration) time.Duration {
	max, _ := time.ParseDuration("0s")
	for _, dur := range durs {
		if dur.Nanoseconds() > max.Nanoseconds() {
			max = dur
		}
	}
	return max
}

// DurationFromProtobuf converts a protobuf duration to a
// `time.Duration`.
// More on protobuf: https://godoc.org/github.com/golang/protobuf/ptypes/duration#Duration
func DurationFromProtobuf(pdur *durationpb.Duration) time.Duration {
	dur, err := time.ParseDuration(
		strconv.Itoa(int((pdur.Seconds*nanosPerSecond)+int64(pdur.Nanos))) + "ns")
	if err != nil {
		panic(err)
	}
	return dur
}

func DurationDaysInt64(dur time.Duration) int64 { return int64(dur.Hours()/24.0) + 1 }

func DurationDays(d time.Duration) float64 { return d.Hours() / 24 }

func DurationYears(d time.Duration) float64 { return DurationDays(d) / 365 }
