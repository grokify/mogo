package timeutil

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

/*

TimeDeltaDow is designed to retrieve a time object x days of week in the past or the future.

// Two Sundays in the future, including today, at 00:00:00
t, err := TimeDeltaDow(time.Now(), time.Sunday, 2, true, true)

// Two Sundays in the future, including today, at present time
t, err := TimeDeltaDow(time.Now(), time.Sunday, 2, true, false)

// Two Sundays ago, not including today, at 00:00:00
t, err := TimeDeltaDow(time.Now(), time.Sunday, -2, false, true)

// Two Sundays ago, not including today, at present time
t, err := TimeDeltaDow(time.Now(), time.Sunday, -2, false, false)

*/

func TimeDeltaDow(base time.Time, wantDow time.Weekday, deltaUnits int, wantInclusive bool, wantStartOfDay bool) (time.Time, error) {
	return TimeDeltaDowInt(base, wantDow, deltaUnits, wantInclusive, wantStartOfDay)
}

func TimeDeltaDowString(base time.Time, wantDowS string, deltaUnits int, wantInclusive bool, wantStartOfDay bool) (time.Time, error) {
	wantDow, err := ParseDayOfWeek(wantDowS)
	if err != nil {
		return time.Now(), err
	}
	return TimeDeltaDowInt(base, wantDow, deltaUnits, wantInclusive, wantStartOfDay)
}

func TimeDeltaDowInt(base time.Time, wantDow time.Weekday, deltaUnits int, wantInclusive bool, wantStartOfDay bool) (time.Time, error) {
	deltaUnitsAbs := deltaUnits
	if deltaUnitsAbs < 1 {
		deltaUnitsAbs *= -1
	}
	deltaDays := int(0)
	if deltaUnits < 0 {
		deltaDays = DaysAgoDow(base.Weekday(), wantDow, wantInclusive)
	} else if deltaUnits > 0 {
		deltaDays = DaysToDow(base.Weekday(), wantDow, wantInclusive)
	} else {
		return base, errors.New("delta units cannot be 0")
	}
	if deltaUnitsAbs > 1 {
		additional := deltaUnitsAbs - 1
		deltaDays += 7 * additional
	}
	if deltaUnits < 0 {
		deltaDays *= -1
	}
	t1 := base.AddDate(0, 0, deltaDays)
	if !wantStartOfDay {
		return t1, nil
	}
	t2 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	return t2, nil
}

func DaysAgoDowStrings(baseDowS string, wantDowS string, wantInclusive bool) (int, error) {
	daysAgo := int(0)
	baseDow, err := ParseDayOfWeek(baseDowS)
	if err != nil {
		return daysAgo, err
	}
	wantDow, err := ParseDayOfWeek(wantDowS)
	if err != nil {
		return daysAgo, err
	}
	return DaysAgoDow(baseDow, wantDow, wantInclusive), nil
}

func DaysAgoDow(baseDow, wantDow time.Weekday, wantInclusive bool) int {
	baseDow = WeekdayNormalized(baseDow)
	wantDow = WeekdayNormalized(wantDow)

	deltaDays1 := int(baseDow) - int(wantDow)
	deltaDays2 := deltaDays1
	if deltaDays2 < 0 {
		deltaDays2 += 7
	}
	if !wantInclusive && deltaDays2 == 0 {
		deltaDays2 = 7
	}
	return deltaDays2
}

func DaysToDowStrings(baseDowS string, wantDowS string, wantInclusive bool) (int, error) {
	daysAgo := int(0)
	baseDow, err := ParseDayOfWeek(baseDowS)
	if err != nil {
		return daysAgo, err
	}
	wantDow, err := ParseDayOfWeek(wantDowS)
	if err != nil {
		return daysAgo, err
	}
	return DaysToDow(baseDow, wantDow, wantInclusive), nil
}

func DaysToDow(baseDow, wantDow time.Weekday, wantInclusive bool) int {
	baseDow = WeekdayNormalized(baseDow)
	wantDow = WeekdayNormalized(wantDow)

	deltaDays1 := int(wantDow) - int(baseDow)
	deltaDays2 := deltaDays1
	if deltaDays2 < 0 {
		deltaDays2 += 7
	}
	if !wantInclusive && deltaDays2 == 0 {
		deltaDays2 = 7
	}
	return deltaDays2
}

func ParseDayOfWeek(value string) (time.Weekday, error) {
	valueLc := strings.ToLower(strings.TrimSpace(value))
	days := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}
	mapping := map[string]int{}
	for i, dow := range days {
		mapping[dow] = i
	}
	if dow, ok := mapping[valueLc]; ok {
		return time.Weekday(dow), nil
	}
	return 0, fmt.Errorf("name of day not found in English: %v", value)
}
