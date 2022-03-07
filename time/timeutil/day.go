package timeutil

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/mogo/type/stringsutil"
)

var days = [...]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

func ParseWeekday(s string) (time.Weekday, error) {
	for i, day := range days {
		if stringsutil.Equal(s, day, true, true) {
			return time.Weekday(i), nil
		}
	}
	return time.Weekday(0), fmt.Errorf("cannot parse weekday: %s", s)
}

// var rxYyyy = regexp.MustCompile(`^[0-9]+$`)

func IntToBaseXString(baseX, val int) string {
	return big.NewInt(int64(val)).Text(baseX)
}

func TimeToDd6(dt time.Time, reverse bool) (string, error) {
	yyyy := dt.Year()
	if yyyy < 0 {
		return "", fmt.Errorf("year is negative [%d]", yyyy)
	} else if yyyy > 9999 {
		return "", fmt.Errorf("year is 5 or more digits [%d]", yyyy)
	}
	dd6 := ""
	if reverse {
		dd6 = big.NewInt(int64(dt.Day())).Text(36) +
			big.NewInt(int64(dt.Month())).Text(36) +
			fmt.Sprintf("%04d", yyyy)
	} else {
		dd6 = fmt.Sprintf("%04d", yyyy) +
			big.NewInt(int64(dt.Month())).Text(36) +
			big.NewInt(int64(dt.Day())).Text(36)
	}
	if len(dd6) != 6 {
		return "", fmt.Errorf("result is not 6 characters [%s]", dd6)
	}
	return dd6, nil
}

var rxDd6Fwd = regexp.MustCompile(`^([0-9]{4})([0-9a-c])([0-9a-v])$`)
var rxDd6Rev = regexp.MustCompile(`^([0-9a-v])([0-9a-c])([0-9]{4})$`)

func Dd6ToTime(dd6 string, reverse bool) (time.Time, error) {
	dd6 = strings.ToLower(strings.TrimSpace(dd6))
	yyyy := ""
	mm36 := ""
	dd36 := ""
	dd6type := ""
	if reverse {
		dd6type = "dd6rev"
		m := rxDd6Rev.FindStringSubmatch(dd6)
		if len(m) == 0 {
			return time.Now(), fmt.Errorf("Dd6Rev Date not parseable [%v]", dd6)
		}
		dd36 = m[1]
		mm36 = m[2]
		yyyy = m[3]
	} else {
		dd6type = "dd6fwd"
		m := rxDd6Fwd.FindStringSubmatch(dd6)
		if len(m) == 0 {
			return time.Now(), fmt.Errorf("Dd6Rev Date not parseable [%v]", dd6)
		}
		yyyy = m[1]
		mm36 = m[2]
		dd36 = m[3]
	}
	y, err := strconv.Atoi(yyyy)
	if err != nil {
		return time.Now(), fmt.Errorf("%s invalid year [%v]", dd6type, dd6)
	}
	m := base36DigitToInt(mm36)
	d := base36DigitToInt(dd36)
	if m < 1 || m > 12 {
		return time.Now(), fmt.Errorf("%s invalid month [%v]", dd6type, dd6)
	}
	if d < 1 || d > 31 {
		return time.Now(), fmt.Errorf("%s invalid day [%v]", dd6type, dd6)
	}
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC), nil
}

const base36Dictionary = "0123456789abcdefghijklmnopqrstuvwxyz"

func base36DigitToInt(val string) int {
	for i, b36 := range base36Dictionary {
		if val == string(b36) {
			return i
		}
	}
	return -1
}
