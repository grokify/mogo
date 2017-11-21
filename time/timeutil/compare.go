// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"time"
)

// IsGreaterThan compares two times and returns true if the left
// time is greater than the right time.
func IsGreaterThan(timeLeft time.Time, timeRight time.Time, orEqual bool) bool {
	durZero, _ := time.ParseDuration("0ns")
	durDelta := timeLeft.Sub(timeRight)
	if durDelta > durZero {
		return true
	} else if orEqual && durDelta == durZero {
		return true
	}
	return false
}

// IsLessThan compares two times and returns true if the left
// time is less than the right time.
func IsLessThan(timeLeft time.Time, timeRight time.Time, orEqual bool) bool {
	durZero, _ := time.ParseDuration("0ns")
	durDelta := timeLeft.Sub(timeRight)
	if durDelta < durZero {
		return true
	} else if orEqual && durDelta == durZero {
		return true
	}
	return false
}

// MinTime returns minTime if time in question is less than min time.
func MinTime(t, min time.Time) time.Time {
	if IsLessThan(t, min, false) {
		return min
	}
	return t
}

// MaxTime returns maxTime if time in question is greater than max time.
func MaxTime(t, max time.Time) time.Time {
	if IsGreaterThan(t, max, false) {
		return max
	}
	return t
}

// GreaterTime returns the greater of two times.
func GreaterTime(t1, t2 time.Time) time.Time {
	if IsGreaterThan(t1, t2, false) {
		return t1
	}
	return t2
}

// LesserTime returns the lesser of two times.
func LesserTime(t1, t2 time.Time) time.Time {
	if IsLessThan(t1, t2, false) {
		return t1
	}
	return t2
}

// MinMax takes two times and returns the earlier time first.
func MinMax(min, max time.Time) (time.Time, time.Time) {
	if IsGreaterThan(min, max, false) {
		return max, min
	}
	return min, max
}
