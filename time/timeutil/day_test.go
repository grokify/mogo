package timeutil

import (
	"testing"
	"time"
)

var dd6Tests = []struct {
	rfc3339 string
	dd6fwd  string
	dd6rev  string
}{
	{"0006-01-02T15:04:00Z", "000612", "210006"},
	{"2006-01-02T15:04:00Z", "200612", "212006"},
	{"2006-10-02T15:04:00Z", "2006a2", "2a2006"},
	{"2006-11-02T15:04:00Z", "2006b2", "2b2006"},
	{"2006-12-02T15:04:00Z", "2006c2", "2c2006"},
	{"2006-12-12T15:04:00Z", "2006cc", "cc2006"},
	{"2006-12-31T15:04:00Z", "2006cv", "vc2006"},
}

// TestDd6 handles dd6 parsing tests
func TestDd6(t *testing.T) {
	for _, tt := range dd6Tests {
		dt, err := time.Parse(time.RFC3339, tt.rfc3339)
		if err != nil {
			t.Errorf("time.Parse(time.RFC3339,tt.rfc3339) Error: with [%v], err [%v]",
				tt.rfc3339, err)
		}
		dd6fwd, err := TimeToDd6(dt, false)
		if err != nil {
			t.Errorf("timeutil.TimeToDd6(time.Time, false) Error: with [%v], err [%v]",
				dt, err)
		}
		if dd6fwd != tt.dd6fwd {
			t.Errorf("timeutil.TimeToDd6(time.Time, false) Mismatch: with [%v], want [%v], got [%v]",
				dt, tt.dd6fwd, dd6fwd)
		}

		dd6rev, err := TimeToDd6(dt, true)
		if err != nil {
			t.Errorf("timeutil.TimeToDd6(time.Time, true) Error: with [%v], err [%v]",
				dt, err)
		}
		if dd6rev != tt.dd6rev {
			t.Errorf("timeutil.TimeToDd6(time.Time, true) Mismatch: with [%v], want [%v], got [%v]",
				dt, tt.dd6rev, dd6rev)
		}

		dtFwd, err := Dd6ToTime(dd6fwd, false)
		if err != nil {
			t.Errorf("timeutil.Dd6ToTime(dd6fwd, false) Error: with [%v], err [%v]",
				dd6fwd, err)
		}
		if dtFwd.Year() != dt.Year() ||
			dtFwd.Month() != dt.Month() ||
			dtFwd.Day() != dt.Day() {
			t.Errorf("timeutil.Dd6ToTime(dd6fwd, false) Error: want [%v], got [%v]",
				dt, dtFwd)
		}

		dtRev, err := Dd6ToTime(dd6rev, true)
		if err != nil {
			t.Errorf("timeutil.Dd6ToTime(dd6rev, true) Error: with [%v], err [%v]",
				dd6rev, err)
		}
		if dtRev.Year() != dt.Year() ||
			dtFwd.Month() != dt.Month() ||
			dtFwd.Day() != dt.Day() {
			t.Errorf("timeutil.Dd6ToTime(dd6rev, true) Error: want [%v], got [%v]",
				dt, dtRev)
		}

	}
}
