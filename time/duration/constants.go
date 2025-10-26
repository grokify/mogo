package duration

import "time"

type Unit string

const (
	UnitYear        Unit = "year"
	UnitMonth       Unit = "month"
	UnitWek         Unit = "week"
	UnitDay         Unit = "day"
	UnitHour        Unit = "hour"
	UnitMinute      Unit = "minute"
	UnitSecond      Unit = "second"
	UnitMillisecond Unit = "millisecond"
	UnitMicrosecond Unit = "microsecond"
	UnitNanosecond  Unit = "nanosecond"

	UnitBusinessDay = "businessday"

	UnitSuffixNanosecond  = "ns"
	UnitSuffixMicrosecond = "us"
	UnitSuffixMillisecond = "ms"
	UnitSuffixSecond      = "s"
	UnitSuffixMinute      = "m"
	UnitSuffixHour        = "h"
	UnitSuffixDay         = "d"
	UnitSuffixWeek        = "w"

	Day         = 24 * time.Hour
	Week        = 7 * Day
	Year        = 365 * Day
	Decade      = 10 * Year
	Score       = 2 * Decade
	Century     = 10 * Decade
	WorkDay     = 8 * time.Hour
	WorkWeek    = 5 * WorkDay
	WorkDay996  = 12 * time.Hour
	WorkWeek996 = 6 * WorkDay996

	HoursPerDay = float32(24)
	DaysPerWeek = float32(7)

	DaySeconds  = 24 * 60 * 60
	WeekSeconds = 7 * DaySeconds
	YearSeconds = (365 * DaySeconds) + (6 * 60 * 60)
)
