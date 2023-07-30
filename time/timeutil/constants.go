package timeutil

import "time"

const (
	SecondsPerYear = (365 * 24 * 60 * 60) + (6 * 60 * 60)
	SecondsPerWeek = 7 * 24 * 60 * 60
	SecondsPerDay  = 24 * 60 * 60

	NanosPerSecond      = int64(1000000000)
	NanosPerMicrosecond = NanosPerSecond / 1000000
	NanosPerMillisecond = NanosPerSecond / 1000
	NanosPerMinute      = NanosPerSecond * 60
	NanosPerHour        = NanosPerMinute * 60
	NanosPerDay         = NanosPerHour * 24
	NanosPerWeek        = NanosPerDay * 7
	NanosPerSecondSub1  = NanosPerSecond - 1

	DurationDay  = 24 * time.Hour
	DurationWeek = 7 * DurationDay

	MonthsEN = `["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]`
)
