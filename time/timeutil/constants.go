package timeutil

import "time"

const (
	Day  = 24 * time.Hour
	Week = 7 * Day

	DaysPerWeek = float32(7)
	HoursPerDay = float32(24)

	SecondsPerDay  = 24 * 60 * 60
	SecondsPerWeek = 7 * SecondsPerDay
	SecondsPerYear = (365 * SecondsPerDay) + (6 * 60 * 60)

	MonthsEN = `["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]`
)
