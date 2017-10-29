// timeutil provides a set of time utilities including comparisons,
// conversion to "DT8" int32 and "DT14" int64 formats and other
// capabilities.
package timeutil

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	DT14            = "20060102150405"
	DT6             = "200601"
	DT8             = "20060102"
	RFC3339Min      = "0000-01-01T00:00:00Z"
	RFC3339Zero     = "0001-01-01T00:00:00Z"
	RFC3339YMD      = "2006-01-02"
	ISO8601YM       = "2006-01"
	ISO8601Z2       = "2006-01-02T15:04:05-07"
	ISO8601Z4       = "2006-01-02T15:04:05-0700"
	ISO8601ZCompact = "20060102T150405Z"
)

var FormatMap = map[string]string{
	"RFC3339":    time.RFC3339,
	"RFC3339YMD": RFC3339YMD,
	"ISO8601YM":  ISO8601YM,
}

func GetFormat(formatName string) (string, error) {
	format, ok := FormatMap[strings.TrimSpace(formatName)]
	if !ok {
		return "", errors.New(fmt.Sprintf("Format Not Found: %v", format))
	}
	return format, nil
}

func FormatQuarter(t time.Time) string {
	return fmt.Sprintf("%d Q%d", t.Year(), MonthToQuarter(int(t.Month())))
}

func TimeRFC3339Zero() time.Time {
	t0, _ := time.Parse(time.RFC3339, RFC3339Zero)
	return t0
}

type RFC3339YMDTime struct {
	time.Time
}

func (t *RFC3339YMDTime) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse(RFC3339YMD, strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

func (t RFC3339YMDTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Time.Format(RFC3339YMD) + `"`), nil
}
