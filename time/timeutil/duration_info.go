package timeutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 00:00:00,309 - 00:00:07,074	 And in conclusion, we have found MySQL to be an excellent database for our website. Any questions?	S1

// DurationInfo is a struct that holds integer values
// for each time unit including hours, minutes, seconds
// milliseconds, microseconds, and nanoseconds.
type DurationInfo struct {
	Hours        int64
	Minutes      int64
	Seconds      int64
	Milliseconds int64
	Microseconds int64
	Nanoseconds  int64
}

// NewDurationInfo returns a DurationInfo struct
// for a duration in nanos.
func NewDurationInfo(dur time.Duration) DurationInfo {
	workingNanos := dur.Nanoseconds()
	dinfo := DurationInfo{}
	if workingNanos >= NanosPerHour {
		hrs := float64(workingNanos) / float64(NanosPerHour)
		hrsInt64 := int64(hrs)
		dinfo.Hours = hrsInt64
		workingNanos = workingNanos - (hrsInt64 * NanosPerHour)
	}
	if workingNanos >= NanosPerMinute {
		min := float64(workingNanos) / float64(NanosPerMinute)
		minInt64 := int64(min)
		dinfo.Minutes = minInt64
		workingNanos = workingNanos - (minInt64 * NanosPerMinute)
	}
	if workingNanos >= NanosPerSecond {
		sec := float64(workingNanos) / float64(NanosPerSecond)
		secInt64 := int64(sec)
		dinfo.Seconds = secInt64
		//workingNanos = workingNanos - (secInt64 * nanosPerSecond)
	}
	return dinfo
}

// ParseDurationInfoStrings returns a DurationInfo object for
// various time units.
func ParseDurationInfoStrings(hr, mn, sc, ms, us, ns string) (DurationInfo, error) {
	dur := DurationInfo{}
	hr = strings.TrimSpace(hr)
	if len(hr) > 0 {
		hours, err := strconv.Atoi(hr)
		if err != nil {
			return dur, err
		}
		dur.Hours = int64(hours)
	}
	if len(mn) > 0 {
		minutes, err := strconv.Atoi(mn)
		if err != nil {
			return dur, err
		}
		dur.Minutes = int64(minutes)
	}
	if len(sc) > 0 {
		seconds, err := strconv.Atoi(sc)
		if err != nil {
			return dur, err
		}
		dur.Seconds = int64(seconds)
	}
	if len(ms) > 0 {
		milliseconds, err := strconv.Atoi(ms)
		if err != nil {
			return dur, err
		}
		dur.Milliseconds = int64(milliseconds)
	}
	if len(us) > 0 {
		microseconds, err := strconv.Atoi(us)
		if err != nil {
			return dur, err
		}
		dur.Microseconds = int64(microseconds)
	}
	if len(ns) > 0 {
		nanoseconds, err := strconv.Atoi(ns)
		if err != nil {
			return dur, err
		}
		dur.Nanoseconds = int64(nanoseconds)
	}
	return dur, nil
}

// TotalNanoseconds returns the total number of nanoseconds
// represented by the duration.
func (di *DurationInfo) TotalNanoseconds() int64 {
	return (di.Hours * NanosPerHour) +
		(di.Minutes * NanosPerMinute) +
		(di.Seconds * NanosPerSecond) +
		(di.Milliseconds * NanosPerMillisecond) +
		(di.Microseconds * NanosPerMicrosecond) +
		di.Nanoseconds
}

// Duration returns a `time.Duration` struct representing
// the duration.
func (di *DurationInfo) Duration() time.Duration {
	dur, err := time.ParseDuration(strconv.Itoa(int(di.TotalNanoseconds())) + "ns")
	if err != nil {
		panic(err)
	}
	return dur
}

// FormatDurationInfoMinSec returns the duration as a simple string
// like 01:01.
func FormatDurationInfoMinSec(di DurationInfo) string {
	min := di.Hours*60 + di.Minutes
	sec := di.Seconds
	return fmt.Sprintf(`%02d:%02d`, min, sec)
}
