package timeutil

import "time"

// InRange checks to see if a time is within a provided time range
// with options whether the begin and end ranges are inclusive or
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
	} else {
		return time.Duration(0)
		panic("E_MORE_THAN_2_RANGES")
	}
	return time.Duration(0)
}

// TimeRange represents a time range with a max and min value.
type TimeRange struct {
	Max     time.Time
	Min     time.Time
	HaveMax bool
	HaveMin bool
}

// Insert updates a time range min and max values for a given time.
func (tr *TimeRange) Insert(t time.Time) {
	tr.InsertMax(t)
	tr.InsertMin(t)
}

// InsertMax updates a time range max value for a given time.
func (tr *TimeRange) InsertMax(t time.Time) {
	if !tr.HaveMax {
		tr.Max = t
		tr.HaveMax = true
	} else if IsGreaterThan(t, tr.Max, false) {
		tr.Max = t
	}
}

// InsertMin updates a time range min value for a given time.
func (tr *TimeRange) InsertMin(t time.Time) {
	if !tr.HaveMin {
		tr.Min = t
		tr.HaveMin = true
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
	if dur.Nanoseconds() > 0 {
		return true
	}
	return false
}
