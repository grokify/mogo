// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

const (
	SecondsPerYear = (365 * 24 * 60 * 60) + (6 * 60 * 60)
	SecondsPerWeek = 7 * 24 * 60 * 60
	SecondsPerDay  = 24 * 60 * 60

	NanosPerSecond      = int64(1000000000)
	NanosPerMicrosecond = NanosPerSecond / 1000000
	NanosPerMillisecond = NanosPerSecond / 1000
	NanosPerMinute      = NanosPerSecond * 60
	NanosPerHour        = NanosPerMinute * 24

	MonthsEN               = `["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]`
	MillisToNanoMultiplier = 1000000
)
