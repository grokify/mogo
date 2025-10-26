package duration

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/grokify/mogo/strconv/strconvutil"
)

// DurationStringUnit converts a Duration to a string with a fixed time unit.
// If an invalid `unit`, is provided, the function will return the output of
// `time.Duration.String()`.
func DurationStringUnit(d, unit time.Duration, prec int, addSuffix bool) string {
	switch unit {
	case time.Nanosecond:
		if v := strconv.Itoa(int(d)); addSuffix {
			return v + UnitSuffixNanosecond
		} else {
			return v
		}
	case time.Microsecond:
		if v := strconv.Itoa(int(d.Microseconds())); addSuffix {
			return v + UnitSuffixMicrosecond
		} else {
			return v
		}
	case time.Millisecond:
		if v := strconv.Itoa(int(d.Milliseconds())); addSuffix {
			return v + UnitSuffixMillisecond
		} else {
			return v
		}
	case time.Second:
		if v := strconvutil.Ftoa(d.Seconds(), prec); addSuffix {
			return v + UnitSuffixSecond
		} else {
			return v
		}
	case time.Minute:
		if v := strconvutil.Ftoa(d.Minutes(), prec); addSuffix {
			return v + UnitSuffixMinute
		} else {
			return v
		}
	case time.Hour:
		if v := strconvutil.Ftoa(d.Hours(), prec); addSuffix {
			return v + UnitSuffixHour
		} else {
			return v
		}
	case Day:
		if v := strconvutil.Ftoa(d.Hours()/24, prec); addSuffix {
			return v + UnitSuffixDay
		} else {
			return v
		}
	case Week:
		if v := strconvutil.Ftoa(d.Hours()/24/7, prec); addSuffix {
			return v + UnitSuffixWeek
		} else {
			return v
		}
	default:
		return d.String()
	}
}

func DurationUnitSuffix(unit time.Duration) string {
	switch unit {
	case time.Nanosecond:
		return UnitSuffixNanosecond
	case time.Microsecond:
		return UnitSuffixMicrosecond
	case time.Millisecond:
		return UnitSuffixMillisecond
	case time.Second:
		return UnitSuffixSecond
	case time.Minute:
		return UnitSuffixMinute
	case time.Hour:
		return UnitSuffixHour
	case Day:
		return UnitSuffixDay
	case Week:
		return UnitSuffixWeek
	default:
		return ""
	}
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
	} else {
		return fmt.Sprintf("%v min %v sec", int(dur.Minutes()), modSeconds), nil
	}
}
