package timeutil

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// InRange checks to see if a time is within a provided time range
// with options whether the start and end ranges are inclusive or
// exclusive. Exclusive ranges are the default.
func InRange(rangeStart, rangeEnd, needle time.Time, incStart, incEnd bool) bool {
	rangeStart, rangeEnd = MinMax(rangeStart, rangeEnd)
	if (incStart && needle.Before(rangeStart)) ||
		(needle.Before(rangeStart) || needle.Equal(rangeStart)) ||
		(incEnd && needle.After(rangeEnd)) ||
		(needle.After(rangeEnd) || needle.Equal(rangeEnd)) {
		return false
	}
	return true
}

type TimeRanges []*TimeRange

func (trs TimeRanges) FilterNonZero() TimeRanges {
	nonzero := TimeRanges{}
	for _, tr := range trs {
		if tr.Duration().Nanoseconds() > 0 {
			nonzero = append(nonzero, tr)
		}
	}
	return nonzero
}

func (trs TimeRanges) IntersectionAny() time.Duration {
	rangesNonZero := trs.FilterNonZero()
	if len(rangesNonZero) == 0 || len(rangesNonZero) == 1 {
		return time.Duration(0)
	} else if len(rangesNonZero) == 2 {
		t1 := rangesNonZero[0]
		t2 := rangesNonZero[1]
		return t1.IntersectionDuration(*t2)
	}
	return time.Duration(0)
}

// TimeRange represents a time range with a max and min value.
type TimeRange struct {
	Max    time.Time
	Min    time.Time
	MinSet bool
	MaxSet bool
}

var rxParseTimeRange = regexp.MustCompile(`^([0-9]+)([MQHY])([0-9]+)$`)

// ParseTimeRangeInterval takes a string in the form of `YYYY[MQY]XX`.
func ParseTimeRangeInterval(s string) (TimeRange, error) {
	s1 := strings.ToUpper(strings.TrimSpace(s))
	m := rxParseTimeRange.FindStringSubmatch(s1)
	if len(m) == 0 {
		return TimeRange{}, fmt.Errorf("cannot parse time range rx (%s)", s)
	}
	yInt, err := strconv.Atoi(m[1])
	if err != nil {
		panic(err)
	}
	intVal, err := strconv.Atoi(m[3])
	if err != nil {
		panic(err)
	}
	// fmtutil.PrintJSON(m)
	switch m[2] {
	case "H":
		if intVal != 1 && intVal != 2 {
			return TimeRange{}, fmt.Errorf("invalid interval (%s)", s)
		}
		switch intVal {
		case 1:
			return TimeRange{
				Min:    time.Date(yInt, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
				Max:    time.Date(yInt, time.Month(7), 1, 0, 0, 0, 0, time.UTC).Add(-1),
				MinSet: true,
				MaxSet: true}, nil
		case 2:
			return TimeRange{
				Min:    time.Date(yInt, time.Month(7), 1, 0, 0, 0, 0, time.UTC),
				Max:    time.Date(yInt+1, time.Month(1), 1, 0, 0, 0, 0, time.UTC).Add(-1),
				MinSet: true,
				MaxSet: true}, nil
		}
	case "Q":
		if intVal < 1 || intVal > 4 {
			return TimeRange{}, fmt.Errorf("invalid interval (%s)", s)
		}
		switch intVal {
		case 1:
			return TimeRange{
				Min:    time.Date(yInt, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
				Max:    time.Date(yInt, time.Month(4), 1, 0, 0, 0, 0, time.UTC).Add(-1),
				MinSet: true,
				MaxSet: true}, nil
		case 2:
			return TimeRange{
				Min:    time.Date(yInt, time.Month(4), 1, 0, 0, 0, 0, time.UTC),
				Max:    time.Date(yInt, time.Month(7), 1, 0, 0, 0, 0, time.UTC).Add(-1),
				MinSet: true,
				MaxSet: true}, nil
		case 3:
			return TimeRange{
				Min:    time.Date(yInt, time.Month(7), 1, 0, 0, 0, 0, time.UTC),
				Max:    time.Date(yInt, time.Month(10), 1, 0, 0, 0, 0, time.UTC).Add(-1),
				MinSet: true,
				MaxSet: true}, nil
		case 4:
			return TimeRange{
				Min:    time.Date(yInt, time.Month(10), 1, 0, 0, 0, 0, time.UTC),
				Max:    time.Date(yInt+1, time.Month(1), 1, 0, 0, 0, 0, time.UTC).Add(-1),
				MinSet: true,
				MaxSet: true}, nil
		}
	}
	return TimeRange{}, fmt.Errorf("time range not supported (%s)", s)
}

func (tr *TimeRange) Contains(t time.Time, inclusiveMin, inclusiveMax bool) (bool, error) {
	if !tr.MinSet || !tr.MaxSet {
		return false, errors.New("timerange must hvae min and max both set")
	}
	if t.Before(tr.Min) || t.After(tr.Max) ||
		(!inclusiveMin && t.Equal(tr.Min)) ||
		(!inclusiveMax && t.Equal(tr.Max)) {
		return false, nil
	} else {
		return true, nil
	}
}

// Insert updates a time range min and max values for a given time.
func (tr *TimeRange) Insert(t time.Time) {
	tr.InsertMax(t)
	tr.InsertMin(t)
}

// InsertMax updates a time range max value for a given time.
func (tr *TimeRange) InsertMax(t time.Time) {
	if !tr.MaxSet {
		tr.Max = t
		tr.MaxSet = true
	} else if IsGreaterThan(t, tr.Max, false) {
		tr.Max = t
	}
}

// InsertMin updates a time range min value for a given time.
func (tr *TimeRange) InsertMin(t time.Time) {
	if !tr.MinSet {
		tr.Min = t
		tr.MinSet = true
	} else if IsLessThan(t, tr.Min, false) {
		tr.Min = t
	}
}

func (tr *TimeRange) Normalize() {
	if tr.Min.After(tr.Max) {
		tmp := tr.Min
		tr.Min = tr.Max
		tr.Max = tmp
	}
}

func (tr *TimeRange) Duration() time.Duration {
	tr.Normalize()
	return tr.Max.Sub(tr.Min)
}

func (tr *TimeRange) Nanoseconds() uint64 {
	tr.Normalize()
	dur := tr.Max.Sub(tr.Min)
	if dur.Nanoseconds() < 0 {
		panic("E_TIMERANGE_DURATION_IS_NEGATIVE")
	}
	return uint64(dur.Nanoseconds())
}

func (tr *TimeRange) IntersectionDuration(tr2 TimeRange) time.Duration {
	tr.Normalize()
	tr2.Normalize()
	if tr2.Min.After(tr.Max) || tr2.Max.Before(tr.Min) {
		// No overlap
		return time.Duration(0)
	} else if (tr2.Min.Equal(tr.Min) || tr2.Min.After(tr.Min)) &&
		(tr2.Max.Equal(tr.Max) || tr2.Max.Before(tr.Max)) {
		// TR2 Completely within TR1
		return tr2.Duration()
	} else if (tr.Min.Equal(tr2.Min) || tr.Min.After(tr2.Min)) &&
		(tr.Max.Equal(tr2.Max) || tr.Max.Before(tr2.Max)) {
		// TR1 Completely within TR2
		return tr.Duration()
	} else if tr.Min.Before(tr2.Min) {
		return tr.Max.Sub(tr2.Min)
	} else {
		return tr2.Max.Sub(tr.Min)
	}
}

func (tr *TimeRange) HasIntersection(tr2 TimeRange) bool {
	dur := tr.IntersectionDuration(tr2)
	return dur.Nanoseconds() > 0
}
