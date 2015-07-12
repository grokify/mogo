package timeutil

import (
	"errors"
	"strings"
	"time"
)

/*

NowDowDeltaStrings is designed to retrieve a time object x days of week in the past or the future.

// Two Sundays in the future, including today, at 00:00:00
t, err := NowDowDeltaStrings("Sunday", 2, true, true)

// Two Sundays in the future, including today, at present time
t, err := NowDowDeltaStrings("Sunday", 2, true, false)

// Two Sundays ago, not including today, at 00:00:00
t, err := NowDowDeltaStrings("Sunday", -2, false, true)

// Two Sundays ago, not including today, at present time
t, err := NowDowDeltaStrings("Sunday", -2, false, false)

*/

func NowDowDeltaString(wantDowS string, deltaUnits int, wantInclusive bool, wantStartOfDay bool) (time.Time, error) {
	now := time.Now()
	deltaUnitsAbs := deltaUnits
	if deltaUnitsAbs < 1 {
		deltaUnitsAbs *= -1
	}
	deltaDays := int(0)
	if deltaUnits < 0 {
		deltaDaysTry, err := DaysAgoDowStrings(now.Weekday().String(), wantDowS, wantInclusive)
		if err != nil {
			return now, err
		}
		deltaDays = deltaDaysTry
	} else if deltaUnits > 0 {
		deltaDaysTry, err := DaysToDowStrings(now.Weekday().String(), wantDowS, wantInclusive)
		if err != nil {
			return now, err
		}
		deltaDays = deltaDaysTry
	}
	if deltaUnitsAbs > 1 {
		additional := deltaUnitsAbs - 1
		deltaDays += 7 * additional
	}
	if deltaUnits < 0 {
		deltaDays *= -1
	}
	t1 := now.AddDate(0, 0, deltaDays)
	if !wantStartOfDay {
		return t1, nil
	}
	t2 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	return t2, nil
}

func DaysAgoDowStrings(baseDowS string, wantDowS string, wantInclusive bool) (int, error) {
	days_ago := int(0)
	baseDow, err := ParseDayOfWeek(baseDowS)
	if err != nil {
		return days_ago, err
	}
	wantDow, err := ParseDayOfWeek(wantDowS)
	if err != nil {
		return days_ago, err
	}
	return DaysAgoDow(baseDow, wantDow, wantInclusive)
}

func DaysAgoDow(baseDow int, wantDow int, wantInclusive bool) (int, error) {
	if baseDow < 0 || baseDow > 6 || wantDow < 0 || wantDow > 6 {
		return int(0), errors.New("Day of week is not in [0-6]")
	}
	deltaDays1 := baseDow - wantDow
	deltaDays2 := deltaDays1
	if deltaDays2 < 0 {
		deltaDays2 = 7 + deltaDays2
	}
	if wantInclusive == false && deltaDays2 == 0 {
		deltaDays2 = 7
	}
	return deltaDays2, nil
}

func DaysToDowStrings(baseDowS string, wantDowS string, wantInclusive bool) (int, error) {
	days_ago := int(0)
	baseDow, err := ParseDayOfWeek(baseDowS)
	if err != nil {
		return days_ago, err
	}
	wantDow, err := ParseDayOfWeek(wantDowS)
	if err != nil {
		return days_ago, err
	}
	return DaysToDow(baseDow, wantDow, wantInclusive)
}

func DaysToDow(baseDow int, wantDow int, wantInclusive bool) (int, error) {
	if baseDow < 0 || baseDow > 6 || wantDow < 0 || wantDow > 6 {
		return int(0), errors.New("Day of week is not in [0-6]")
	}
	deltaDays1 := wantDow - baseDow
	deltaDays2 := deltaDays1
	if deltaDays2 < 0 {
		deltaDays2 = 7 + deltaDays2
	}
	if wantInclusive == false && deltaDays2 == 0 {
		deltaDays2 = 7
	}
	return deltaDays2, nil
}

func ParseDayOfWeek(value string) (int, error) {
	value = strings.ToLower(value)
	days := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}
	mapping := map[string]int{}
	for i, dow := range days {
		mapping[dow] = int(i)
	}
	if dow, ok := mapping[value]; ok {
		return dow, nil
	}
	return -1, errors.New("English name of day not found")
}
