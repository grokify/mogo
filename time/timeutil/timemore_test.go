package timeutil

import (
	"testing"
	"time"
)

var timemoreTests = []struct {
	t            string
	dow          time.Weekday
	yearStart    string
	quarterStart string
	monthStart   string
	weekStart    string
	dayStart     string
	yearHalf     string
	yearQuarter  string
	yearMonth    string
}{
	{"2023-05-15T12:30:30Z", time.Monday, "2023-01-01T00:00:00Z", "2023-04-01T00:00:00Z", "2023-05-01T00:00:00Z", "2023-05-15T00:00:00Z", "2023-05-15T00:00:00Z",
		"2023H1", "2023Q2", "2023M5"},
	{"2023-05-15T12:30:30Z", 7001, "2023-01-01T00:00:00Z", "2023-04-01T00:00:00Z", "2023-05-01T00:00:00Z", "2023-05-15T00:00:00Z", "2023-05-15T00:00:00Z",
		"2023H1", "2023Q2", "2023M5"},
	{"2023-05-15T12:30:30Z", time.Sunday, "2023-01-01T00:00:00Z", "2023-04-01T00:00:00Z", "2023-05-01T00:00:00Z", "2023-05-14T00:00:00Z", "2023-05-15T00:00:00Z",
		"2023H1", "2023Q2", "2023M5"},
	{"2023-05-15T12:30:30Z", 14000, "2023-01-01T00:00:00Z", "2023-04-01T00:00:00Z", "2023-05-01T00:00:00Z", "2023-05-14T00:00:00Z", "2023-05-15T00:00:00Z",
		"2023H1", "2023Q2", "2023M5"},
	{"2023-05-15T12:30:30-08:00", time.Monday, "2023-01-01T00:00:00-08:00", "2023-04-01T00:00:00-08:00", "2023-05-01T00:00:00-08:00", "2023-05-15T00:00:00-08:00", "2023-05-15T00:00:00-08:00", "2023H1", "2023Q2", "2023M5"},
	{"2023-05-15T12:30:30-08:00", time.Tuesday, "2023-01-01T00:00:00-08:00", "2023-04-01T00:00:00-08:00", "2023-05-01T00:00:00-08:00", "2023-05-09T00:00:00-08:00", "2023-05-15T00:00:00-08:00", "2023H1", "2023Q2", "2023M5"},
}

func TestTimeMore(t *testing.T) {
	for _, tt := range timemoreTests {
		dt := MustParse(time.RFC3339, tt.t)

		ys2 := NewTimeMore(dt, tt.dow).YearStart()
		if ys2.Format(time.RFC3339) != tt.yearStart {
			t.Errorf("mismatch TimeMore.YearStart() TimeMore(\"%s\", %v): want [%v], got [%v]", tt.t, tt.dow.String(), tt.yearStart, ys2.Format(time.RFC3339))
		}

		tm := NewTimeMore(dt, tt.dow)

		ys := tm.YearStart()
		if ys.Format(time.RFC3339) != tt.yearStart {
			t.Errorf("mismatch TimeMore.YearStart() TimeMore(\"%s\", %v): want [%v], got [%v]", tt.t, tt.dow.String(), tt.yearStart, ys.Format(time.RFC3339))
		}

		qs := tm.QuarterStart()
		if qs.Format(time.RFC3339) != tt.quarterStart {
			t.Errorf("mismatch TimeMore.QuarterStart() TimeMore(\"%s\", %v): want [%v], got [%v]", tt.t, tt.dow.String(), tt.quarterStart, qs.Format(time.RFC3339))
		}

		ms := tm.MonthStart()
		if ms.Format(time.RFC3339) != tt.monthStart {
			t.Errorf("mismatch TimeMore.MonthStart() TimeMore(\"%s\", %v): want [%v], got [%v]", tt.t, tt.dow.String(), tt.monthStart, ms.Format(time.RFC3339))
		}

		ws := tm.WeekStart()
		if ws.Format(time.RFC3339) != tt.weekStart {
			t.Errorf("mismatch TimeMore.WeekStart() TimeMore(\"%s\", %v): want [%v], got [%v]", tt.t, tt.dow.String(), tt.weekStart, ws.Format(time.RFC3339))
		}

		ds := tm.DayStart()
		if ds.Format(time.RFC3339) != tt.dayStart {
			t.Errorf("mismatch TimeMore.DayStart() TimeMore(\"%s\", %v): want [%v], got [%v]", tt.t, tt.dow.String(), tt.dayStart, ds.Format(time.RFC3339))
		}

		yh := tm.YearHalf()
		if yh != tt.yearHalf {
			t.Errorf("mismatch TimeMore.YearHalf() TimeMore(\"%s\", %v): want [%v], got [%v]", tt.t, tt.dow.String(), tt.yearHalf, yh)
		}

		yq := tm.YearQuarterString()
		if yq != tt.yearQuarter {
			t.Errorf("mismatch TimeMore.YearQuarrter() TimeMore(\"%s\", %v): want [%v], got [%v]", tt.t, tt.dow.String(), tt.yearQuarter, yq)
		}

		ym := tm.YearMonth()
		if ym != tt.yearMonth {
			t.Errorf("mismatch TimeMore.YearMonth() TimeMore(\"%s\", %v): want [%v], got [%v]", tt.t, tt.dow.String(), tt.yearMonth, ym)
		}
	}
}
