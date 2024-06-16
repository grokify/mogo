package timeutil

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NewDuration(day, hour, min, sec, nsec int) time.Duration {
	return time.Duration(int64(day)*int64(time.Hour)*24 +
		int64(hour)*int64(time.Hour) +
		int64(min)*int64(time.Minute) +
		int64(sec)*int64(time.Second) +
		int64(nsec))
}

func NewDurationFloat(day, hour, min, sec float64, nsec int64) time.Duration {
	nps := float64(time.Second)
	return time.Duration(int64(float64(nsec) +
		sec*nps +
		min*60*nps +
		hour*60*60*nps +
		day*24*60*60*nps))
}

func NewDurationStrings(hour, min, sec string) (time.Duration, error) {
	return time.ParseDuration(fmt.Sprintf("%vh%vm%vs", hour, min, sec))
}

// ParseDuration adds days (d), weeks (w), years (y).
func ParseDuration(s string) (time.Duration, error) {
	rx := regexp.MustCompile(`(?i)^\s*(-?\d+)(d|w|y)\s*$`)
	rs := rx.FindStringSubmatch(s)

	if len(rs) > 0 {
		zeroDuration, _ := time.ParseDuration("0s")
		quantity := rs[1]
		units := strings.ToLower(rs[2])
		if i, err := strconv.Atoi(quantity); err != nil {
			return zeroDuration, err
		} else if units == "d" {
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
	if dur, err := time.ParseDuration(s); err != nil {
		panic(err)
	} else {
		return dur
	}
}

func NowDeltaDuration(d time.Duration) time.Time {
	return time.Now().Add(d)
}

func NowDeltaParseDuration(s string) (time.Time, error) {
	if d, err := ParseDuration(s); err != nil {
		return time.Now(), err
	} else {
		return time.Now().Add(d), nil
	}
}

/*
// DurationForNowSubDt8 returns a duartion struct between a Dt8 value and the current time.
func DurationForNowSubDT8(dt8 int32) (time.Duration, error) {
	t, err := TimeForDT8(dt8)
	if err != nil {
		var d time.Duration
		return d, err
	}
	now := time.Now()
	return now.Sub(t), nil
}
*/

// QuarterDuration returns a `time.Duration` representing the calendar quarter for the time provided.
func (tm TimeMore) QuarterDuration() time.Duration {
	start := quarterStart(tm.Time())
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

// SubDuration subtracts one duration from another and returns the result as a `time.Duration`.
func SubDuration(dur1, dur2 time.Duration) time.Duration {
	ns := dur1.Nanoseconds() - dur2.Nanoseconds()
	diff, err := time.ParseDuration(fmt.Sprintf("%dns", ns))
	if err != nil {
		panic("err")
	}
	return diff
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

/*
// DurationFromProtobuf converts a protobuf duration to a
// `time.Duration`.
// More on protobuf: https://godoc.org/github.com/golang/protobuf/ptypes/duration#Duration
func DurationFromProtobuf(pdur *durationpb.Duration) time.Duration {
	dur, err := time.ParseDuration(
		strconv.Itoa(int((pdur.Seconds*NanosPerSecond)+int64(pdur.Nanos))) + "ns")
	if err != nil {
		panic(err)
	}
	return dur
}
*/

func DurationDaysInt64(dur time.Duration) int64 { return int64(dur.Hours()/24.0) + 1 }

func DurationDays(d time.Duration) float64 { return d.Hours() / 24 }

func DurationWeeks(d time.Duration) float64 { return DurationDays(d) / 7 }

func DurationYears(d time.Duration) float64 { return DurationDays(d) / 365 }
