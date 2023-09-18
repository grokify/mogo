package timeutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DurationInfo is a struct that holds integer values
// for each time unit including hours, minutes, seconds
// milliseconds, microseconds, and nanoseconds.
type DurationInfo struct {
	DaysPerWeek  float32
	HoursPerDay  float32
	Weeks        int64
	Days         int64
	Hours        int64
	Minutes      int64
	Seconds      int64
	Milliseconds int64
	Microseconds int64
	Nanoseconds  int64
}

// NewDurationInfo returns a DurationInfo struct for a duration in nanos. If 'daysPerWeek` or
// `hoursPerDay` are set to zero, the default values of 7 and 24 are used.
func NewDurationInfo(d time.Duration, daysPerWeek, hoursPerDay float32) DurationInfo {
	dinfo := DurationInfo{
		DaysPerWeek: daysPerWeek,
		HoursPerDay: hoursPerDay}
	workingNanos := d.Nanoseconds()
	nanosPerWeek := NanosPerWeek
	if daysPerWeek != 0 || hoursPerDay != 0 {
		if daysPerWeek == 0 {
			daysPerWeek = 7
		}
		if hoursPerDay == 0 {
			hoursPerDay = 24
		}
		nanosPerWeek = int64(daysPerWeek * hoursPerDay * float32(NanosPerHour))
	}
	if workingNanos >= nanosPerWeek {
		weeks := float64(workingNanos) / float64(NanosPerHour)
		weeksInt64 := int64(weeks)
		dinfo.Weeks = weeksInt64
		workingNanos = workingNanos - (weeksInt64 * nanosPerWeek)
	}
	nanosPerDay := NanosPerDay
	if hoursPerDay != 0 {
		nanosPerDay = int64(hoursPerDay * float32(NanosPerHour))
	}
	if workingNanos >= nanosPerDay {
		days := float64(workingNanos) / float64(nanosPerDay)
		daysInt64 := int64(days)
		dinfo.Days = daysInt64
		workingNanos = workingNanos - (daysInt64 * nanosPerDay)
	}
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

// Duration returns a `time.Duration` struct. Params for `hoursPerDay` and `daysPerWeek` are
// used for atlernate values such as working hours per day and working days per week, e.g.
// 8 hours per day and 5 days per week.
func (di DurationInfo) Duration(hoursPerDay, daysPerWeek float32) time.Duration {
	dur := time.Duration(di.Nanoseconds) +
		time.Duration(di.Microseconds)*time.Microsecond +
		time.Duration(di.Milliseconds)*time.Millisecond +
		time.Duration(di.Seconds)*time.Second +
		time.Duration(di.Minutes)*time.Minute +
		time.Duration(di.Hours)*time.Hour
	if di.Days != 0 {
		if hoursPerDay != 0 {
			dur += time.Duration(di.Days) * time.Duration(hoursPerDay) * time.Hour
		} else {
			dur += time.Duration(di.Days) * DurationDay
		}
	}
	if di.Weeks != 0 {
		if daysPerWeek != 0 {
			daysPerWeek := time.Duration(daysPerWeek)
			if hoursPerDay != 0 {
				dur += time.Duration(di.Weeks) *
					daysPerWeek *
					time.Duration(hoursPerDay) *
					time.Hour
			} else {
				dur += time.Duration(di.Weeks) *
					daysPerWeek *
					DurationDay
			}
		} else {
			dur += time.Duration(di.Weeks) * DurationWeek
		}
	}
	return dur
}

// ParseDurationInfo converts a Jira human readable string into a `DurationInfo` struct.
func ParseDurationInfo(s string) (DurationInfo, error) {
	// MySQL format: 00:00:00,309 - 00:00:07,074 And in conclusion, we have found MySQL to be an excellent database for our website. Any questions?	S1
	parts := strings.Split(strings.ToLower(s), ",")
	di := DurationInfo{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if len(p) == 0 {
			continue
		}
		ps := strings.Fields(p)
		if len(ps) != 2 {
			return di, fmt.Errorf("cannot parse (%s)", p)
		}
		v, err := strconv.Atoi(ps[0])
		if err != nil {
			return di, err
		}
		v64 := int64(v)
		switch ps[1] {
		case "week", "weeks", "w":
			di.Weeks = v64
		case "day", "days", "d":
			di.Days = v64
		case "hour", "hours", "h":
			di.Hours = v64
		case "minute", "minutes", "m":
			di.Minutes = v64
		case "second", "seconds", "s":
			di.Seconds = v64
		case "millisecond", "milliseconds", "ms":
			di.Milliseconds = v64
		case "microsecond", "microseconds", "us", "µs":
			di.Microseconds = v64
		case "nanosecond", "nanoseconds", "ns":
			di.Nanoseconds = v64
		default:
			return di, fmt.Errorf("cannot parse (%s)", p)
		}
	}
	return di, nil
}

// FormatDurationInfoMinSec returns the duration as a simple string like 01:01.
func FormatDurationInfoMinSec(di DurationInfo) string {
	min := di.Hours*60 + di.Minutes
	sec := di.Seconds
	return fmt.Sprintf(`%02d:%02d`, min, sec)
}

// DurationInfoString represets a set of time duration data. It is useful for converting
// parsed time data into a `time.Duration` struct. `DaysPerWeek` and `HoursPerDay` are provided
// as overrides to standard value of 7 and 24 in the case of business context, e.g. 5 days
// per week and 8 hours per day.
type DurationInfoString struct {
	DaysPerWeek  float32
	HoursPerDay  float32
	Weeks        string
	Days         string
	Hours        string
	Minutes      string
	Seconds      string
	Milliseconds string
	Microseconds string
	Nanoseconds  string
}

func (dis DurationInfoString) Duration() (time.Duration, error) {
	daysPerWeek := dis.DaysPerWeek
	if daysPerWeek <= 0 {
		daysPerWeek = DaysPerWeek
	}
	hoursPerDay := dis.HoursPerDay
	if hoursPerDay <= 0 {
		hoursPerDay = HoursPerDay
	}
	nanosPerDay := int64(float64(NanosPerHour) * float64(hoursPerDay))
	nanosPerWeek := int64(float64(NanosPerHour) * float64(hoursPerDay) * float64(daysPerWeek))

	timeUnitData := []struct {
		v string
		n int64
	}{
		{dis.Weeks, nanosPerWeek},
		{dis.Days, nanosPerDay},
		{dis.Hours, NanosPerHour},
		{dis.Minutes, NanosPerMinute},
		{dis.Seconds, NanosPerSecond},
		{dis.Milliseconds, NanosPerMillisecond},
		{dis.Microseconds, NanosPerMicrosecond},
		{dis.Nanoseconds, 1},
	}

	var nanos int64
	for _, tu := range timeUnitData {
		v := strings.TrimSpace(tu.v)
		if v == "" || v == "0" {
			continue
		} else if f, err := strconv.ParseFloat(v, 64); err != nil {
			return 0, err
		} else {
			nanos += int64(f * float64(tu.n))
		}
	}

	return time.Duration(nanos), nil
}
