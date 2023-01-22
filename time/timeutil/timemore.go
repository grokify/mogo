package timeutil

import (
	"time"
)

// TimeMore is a struct for holding various times related
// to a current time, including year start, quarter start,
// month start, and week start.
type TimeMore struct {
	Time time.Time
	DOW  time.Weekday
}

func NewTimeMore(t time.Time, dow time.Weekday) TimeMore {
	return TimeMore{
		Time: t,
		DOW:  dow}
}

func NewTimeMoreQuarterStartString(yyyyqStr string, dow time.Weekday) (TimeMore, error) {
	dt, err := QuarterStringStartTime(yyyyqStr)
	if err != nil {
		return TimeMore{}, err
	}
	return NewTimeMore(dt, dow), nil
}

func NewTimeMoreQuarterEndString(yyyyqStr string, dow time.Weekday) (TimeMore, error) {
	dt, err := QuarterStringEndTime(yyyyqStr)
	if err != nil {
		return TimeMore{}, err
	}
	return NewTimeMore(dt, dow), nil
}

// func (tm *TimeMore) Time() time.Time         { return tm.thisTime }
// func (tm *TimeMore) DOW() time.Weekday       { return tm.dow }

func (tm TimeMore) DayStart() time.Time     { return dayStart(tm.Time) }
func (tm TimeMore) MonthStart() time.Time   { return monthStart(tm.Time) }
func (tm TimeMore) IsMonthStart() bool      { return isMonthStart(tm.Time) }
func (tm TimeMore) Quarter() int32          { return QuarterInt32ForTime(tm.Time) }
func (tm TimeMore) QuarterStart() time.Time { return quarterStart(tm.Time) }
func (tm TimeMore) QuarterEnd() time.Time   { return quarterEnd(tm.Time) }
func (tm TimeMore) IsQuarterStart() bool    { return isQuarterStart(tm.Time) }
func (tm TimeMore) Year() int32             { return QuarterInt32ToYear(tm.Quarter()) }
func (tm TimeMore) YearStart() time.Time    { return yearStart(tm.Time) }
func (tm TimeMore) YearEnd() time.Time      { return yearEnd(tm.Time) }
func (tm TimeMore) IsYearStart() bool       { return isYearStart(tm.Time) }
func (tm TimeMore) IntervalStart(interval Interval) (time.Time, error) {
	return intervalStart(tm.Time, interval, tm.DOW)
}

// WeekStart takes a `time.Time` struct and a week start day in the `time.Weekday` format.
func (tm TimeMore) WeekStart() time.Time {
	week, err := weekStart(tm.Time, tm.DOW)
	if err != nil {
		panic(err)
	}
	return week
}

// TimeMeta is a struct for holding various times related
// to a current time, including year start, quarter start,
// month start, and week start.
type TimeMeta struct {
	This         time.Time
	YearStart    time.Time
	QuarterStart time.Time
	MonthStart   time.Time
	WeekStart    time.Time
	DayStart     time.Time
}

// TimeMeta returns a TimeMeta struct given `time.Time`
// and `time.Weekday` parameters.
func (tm TimeMore) TimeMeta() TimeMeta {
	// dt = dt.UTC()
	// tm := NewTimeMore(t, 0)
	return TimeMeta{
		This:         tm.Time,
		YearStart:    tm.YearStart(),
		QuarterStart: tm.QuarterStart(),
		MonthStart:   tm.MonthStart(),
		WeekStart:    tm.WeekStart(),
		DayStart:     tm.DayStart(),
	}
	/*
		week, err := tm.WeekStart()
		if err != nil {
			return meta, err
		}
		meta.WeekStart = week
		/
		return meta, nil
	*/
}
