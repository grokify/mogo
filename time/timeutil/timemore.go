package timeutil

import (
	"time"
)

// TimeMore is a struct for holding various times related
// to a current time, including year start, quarter start,
// month start, and week start.
type TimeMore struct {
	thisTime time.Time
	dow      time.Weekday
}

func NewTimeMore(thisTime time.Time, dow time.Weekday) TimeMore {
	return TimeMore{
		thisTime: thisTime,
		dow:      dow}
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

func (tm *TimeMore) Time() time.Time         { return tm.thisTime }
func (tm *TimeMore) DOW() time.Weekday       { return tm.dow }
func (tm *TimeMore) MonthStart() time.Time   { return MonthStart(tm.thisTime) }
func (tm *TimeMore) Quarter() int32          { return QuarterInt32ForTime(tm.thisTime) }
func (tm *TimeMore) QuarterStart() time.Time { return QuarterStart(tm.thisTime) }
func (tm *TimeMore) QuarterEnd() time.Time   { return QuarterEnd(tm.thisTime) }
func (tm *TimeMore) Year() int32             { return QuarterInt32ToYear(tm.Quarter()) }
func (tm *TimeMore) YearStart() time.Time    { return YearStart(tm.thisTime) }
func (tm *TimeMore) YearEnd() time.Time      { return YearEnd(tm.thisTime) }

func (tm *TimeMore) WeekStart() time.Time {
	week, err := WeekStart(tm.thisTime, tm.dow)
	if err != nil {
		panic(err)
	}
	return week
}
