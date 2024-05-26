package customtime

import (
	"fmt"
	"strings"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

type TimeRFC3339 struct {
	time.Time
}

func (ct *TimeRFC3339) UnmarshalJSON(b []byte) (err error) {
	if s := strings.Trim(string(b), "\""); s == "null" {
		ct.Time = time.Time{}
	} else {
		ct.Time, err = time.Parse(timeutil.RFC3339FullDate, s)
	}
	return
}

func (ct *TimeRFC3339) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	} else {
		return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(timeutil.RFC3339FullDate))), nil
	}
}
