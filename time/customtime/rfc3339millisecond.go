package customtime

import (
	"fmt"
	"strings"
	"time"

	"github.com/grokify/simplego/time/timeutil"
)

// TimeRFC3339Milli represents a time that is represened
// to the millisecond (3 sigificant digits vs. Go's 6
// significant digits. This iss useful for some Java
// deployments. Read more here:
// https://stackoverflow.com/a/41678233/1908967
type TimeRFC3339Milli struct {
	time.Time
}

func (ct *TimeRFC3339Milli) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(timeutil.RFC3339Milli, s)
	return
}

func (ct *TimeRFC3339Milli) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(timeutil.RFC3339Milli))), nil
}

var nilTime = (time.Time{}).UnixNano()

func (ct *TimeRFC3339Milli) IsSet() bool {
	return ct.UnixNano() != nilTime
}
