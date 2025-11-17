package markdown

import (
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

func GeneratedAtLocationStringsHuman(dt time.Time, locNames []string) (string, error) {
	return GeneratedAtLocationStrings(timeutil.HumanDateTime, dt, locNames)
}

func GeneratedAtLocationStrings(layout string, dt time.Time, locNames []string) (string, error) {
	locs := []*time.Location{}
	for _, locStr := range locNames {
		if loc, err := time.LoadLocation(locStr); err != nil {
			return "", err
		} else if loc != nil {
			locs = append(locs, loc)
		}
	}
	return GeneratedAt(layout, dt, locs), nil
}

func GeneratedAt(layout string, dt time.Time, locs []*time.Location) string {
	prefix := "Generated at"
	if len(locs) == 0 {
		return prefix + " " + dt.Format(layout)
	} else if len(locs) == 1 {
		if locs[0] != nil {
			dt = dt.In(locs[0])
		}
		return prefix + " " + dt.Format(layout)
	}
	str := prefix + "\n"
	for _, loc := range locs {
		if loc != nil {
			dt = dt.In(loc)
		}
		str += "* " + dt.Format(layout) + "\n"
	}
	return str
}
