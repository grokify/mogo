package timeutil

import (
	"encoding/json"
	"time"
)

func MonthEndDay(year int, month time.Month) int {
	if month == time.December {
		month = time.January
		year++
	} else {
		month++
	}
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1).Day()
}

func MonthNames() []string {
	data := []string{}
	err := json.Unmarshal([]byte(MonthsEN), &data)
	if err != nil {
		panic(err)
	}
	return data
}
