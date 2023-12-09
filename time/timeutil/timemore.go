package timeutil

import (
	"fmt"
	"time"
)

// TimeMore is a struct for holding various times related
// to a current time, including year start, quarter start,
// month start, and week start.
type TimeMore struct {
	time         time.Time
	weekStartDay time.Weekday
}

func NewTimeMore(t time.Time, d time.Weekday) TimeMore {
	return TimeMore{
		time:         t,
		weekStartDay: WeekdayNormalized(d)}
}

func NewTimeMoreQuarterStartString(yyyyq string, d time.Weekday) (TimeMore, error) {
	dt, err := QuarterStringStartTime(yyyyq)
	if err != nil {
		return TimeMore{}, err
	}
	return NewTimeMore(dt, d), nil
}

func NewTimeMoreQuarterEndString(yyyyq string, d time.Weekday) (TimeMore, error) {
	dt, err := QuarterStringEndTime(yyyyq)
	if err != nil {
		return TimeMore{}, err
	}
	return NewTimeMore(dt, d), nil
}

func (tm TimeMore) Time() time.Time            { return tm.time }
func (tm TimeMore) WeekStartDay() time.Weekday { return tm.weekStartDay }

func (tm TimeMore) DayStart() time.Time  { return dayStart(tm.time) }
func (tm TimeMore) DayEnd() time.Time    { return dayEnd(tm.time) }
func (tm TimeMore) WeekStart() time.Time { return weekStart(tm.time, tm.weekStartDay) }
func (tm TimeMore) WeekEnd() time.Time   { return weekEnd(tm.time, tm.weekStartDay) }
func (tm TimeMore) HalfYear() uint8 {
	if tm.time.Month() < time.July {
		return 1
	} else {
		return 2
	}
}
func (tm TimeMore) MonthStart() time.Time   { return monthStart(tm.time) }
func (tm TimeMore) MonthEnd() time.Time     { return monthEnd(tm.time) }
func (tm TimeMore) IsMonthStart() bool      { return isMonthStart(tm.time) }
func (tm TimeMore) Quarter() int32          { return QuarterInt32ForTime(tm.time) }
func (tm TimeMore) QuarterStart() time.Time { return quarterStart(tm.time) }
func (tm TimeMore) QuarterEnd() time.Time   { return quarterEnd(tm.time) }
func (tm TimeMore) IsQuarterStart() bool    { return isQuarterStart(tm.time) }
func (tm TimeMore) Year() int32             { return QuarterInt32ToYear(tm.Quarter()) }
func (tm TimeMore) YearStart() time.Time    { return yearStart(tm.time) }
func (tm TimeMore) YearEnd() time.Time      { return yearEnd(tm.time) }
func (tm TimeMore) IsYearStart() bool       { return isYearStart(tm.time) }
func (tm TimeMore) IntervalStart(interval Interval) (time.Time, error) {
	return intervalStart(tm.time, interval, tm.weekStartDay)
}

// YearHalf returns a string in the format of "2006H1"
func (tm TimeMore) YearHalf() string { return fmt.Sprintf("%dH%d", tm.time.Year(), int(tm.HalfYear())) }

// YearQuarter returns a string in the format of "2006Q1"
func (tm TimeMore) YearQuarter() string {
	return fmt.Sprintf("%dQ%d", tm.time.Year(), int(tm.QuarterCalendar()))
}

// YearMonth returns a string in the format of "2006M1"
func (tm TimeMore) YearMonth() string {
	return fmt.Sprintf("%dM%d", tm.time.Year(), int(tm.time.Month()))
}

// QuarterCalendar returns the quarter of the year specified by tm.Time.
func (tm TimeMore) QuarterCalendar() Yearquarter {
	m := tm.time.Month()
	if m < 4 {
		return Winter // 1
	} else if m < 7 {
		return Spring
	} else if m < 10 {
		return Summer
	} else {
		return Autumn
	}
}

func (tm TimeMore) SeasonMeteorological() Yearquarter {
	m := tm.time.Month()
	if m < 3 {
		return Winter // Starts December 1
	} else if m < 6 {
		return Spring // Starts March 1
	} else if m < 9 {
		return Summer // Starts June 1
	} else if m < 12 {
		return Autumn // Starts September 1
	} else {
		return Winter // Starts December 1
	}
}

// TimeMeta is a struct for holding various times related
// to a current time, including year start, quarter start,
// month start, and week start.
type TimeMeta struct {
	This           time.Time
	WeekStartDay   string
	YearStart      time.Time
	YearEnd        time.Time
	QuarterStart   time.Time
	QuarterEnd     time.Time
	MonthStart     time.Time
	MonthEnd       time.Time
	WeekStart      time.Time
	WeekEnd        time.Time
	DayStart       time.Time
	DayEnd         time.Time
	IsYearStart    bool
	IsQuarterStart bool
	IsMonthStart   bool
}

// TimeMeta returns a TimeMeta struct given `time.Time` and `time.Weekday` parameters.
func (tm TimeMore) TimeMeta() TimeMeta {
	return TimeMeta{
		This:           tm.Time(),
		WeekStartDay:   tm.weekStartDay.String(),
		YearStart:      tm.YearStart(),
		YearEnd:        tm.YearEnd(),
		QuarterStart:   tm.QuarterStart(),
		QuarterEnd:     tm.QuarterEnd(),
		MonthStart:     tm.MonthStart(),
		MonthEnd:       tm.MonthEnd(),
		WeekStart:      tm.WeekStart(),
		WeekEnd:        tm.WeekEnd(),
		DayStart:       tm.DayStart(),
		DayEnd:         tm.DayEnd(),
		IsYearStart:    tm.IsYearStart(),
		IsQuarterStart: tm.IsQuarterStart(),
		IsMonthStart:   tm.IsMonthStart(),
	}
}
