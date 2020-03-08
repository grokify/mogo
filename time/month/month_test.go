package month

import (
	"testing"
	"time"

	"github.com/grokify/gotilla/time/timeutil"
)

var dayofmonthToEnglishTests = []struct {
	v    uint16
	want string
}{
	{uint16(1), "first"},
	{uint16(2), "second"},
	{uint16(3), "third"},
	{uint16(19), "nineteenth"},
	{uint16(20), "twentieth"},
	{uint16(21), "twenty first"},
	{uint16(30), "thirtieth"},
	{uint16(31), "thirty first"},
}

func TestDayofmonthToEnglish(t *testing.T) {
	for _, tt := range dayofmonthToEnglishTests {
		got := DayofmonthToEnglish(tt.v)
		if got != tt.want {
			t.Errorf("DayofmonthToEnglish(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}

var yearMonthBase36Tests = []struct {
	year  uint64
	month uint64
	want  string
}{
	{uint64(1), uint64(1), "002T"},
	{uint64(500), uint64(3), "12KZ"},
	{uint64(1000), uint64(6), "255Y"},
	{uint64(2010), uint64(8), "4B3K"},
	{uint64(2019), uint64(8), "4BSK"},
	{uint64(2038), uint64(1), "4D95"},
	{uint64(9999), uint64(12), "LFJC"},
}

func TestYearMonthBase36(t *testing.T) {
	for _, tt := range yearMonthBase36Tests {
		got := YearMonthBase36(tt.year, tt.month)
		if got != tt.want {
			t.Errorf("YearMonthBase36(%v, %v): want %v, got %v", tt.year, tt.month, tt.want, got)
		}
	}
}

var monthContinuousTests = []struct {
	year    uint64
	month   uint64
	monthc  uint64
	rfc3339 string
}{
	{uint64(0), uint64(1), uint64(1), "0000-01-01T00:00:00Z"},
	{uint64(0), uint64(12), uint64(12), "0000-12-01T00:00:00Z"},
	{uint64(1), uint64(1), uint64(13), "0001-01-01T00:00:00Z"},
	{uint64(1), uint64(12), uint64(24), "0001-12-01T00:00:00Z"},
	{uint64(2), uint64(1), uint64(25), "0002-01-01T00:00:00Z"},
	{uint64(3), uint64(12), uint64(48), "0003-12-01T00:00:00Z"},
	{uint64(4), uint64(1), uint64(49), "0004-01-01T00:00:00Z"},
}

func TestMonthContinuous(t *testing.T) {
	for _, tt := range monthContinuousTests {
		gotMonthC := YearMonthToMonthContinuous(tt.year, tt.month)
		if gotMonthC != tt.monthc {
			t.Errorf("YearMonthToMonthContinuous(%v, %v): want [%v], got [%v]",
				tt.year, tt.month, tt.monthc, gotMonthC)
		}
		gotYear, gotMonth := MonthContinuousToYearMonth(tt.monthc)
		if gotYear != tt.year || gotMonth != tt.month {
			t.Errorf("MonthContinuousToYearMonth(%v): want [%v,%v], got [%v,%v]",
				tt.monthc,
				tt.year, tt.month,
				gotYear, gotMonth)
		}
		dt := MonthContinuousToTime(tt.monthc)
		t3339 := dt.Format(time.RFC3339)
		if t3339 != tt.rfc3339 {
			t.Errorf("MonthContinuousToTime(%v, %v): want [%v], got [%v]",
				tt.year, tt.month, tt.rfc3339, t3339)
		}
		wantDt, err := time.Parse(time.RFC3339, tt.rfc3339)
		if err != nil {
			t.Errorf("TimeToMonthContinuous time.Parse(time.RFC3339, \"%s\") error [%v]",
				tt.rfc3339, err.Error())
		}
		gotMonthc := TimeToMonthContinuous(wantDt)
		if gotMonthc != tt.monthc {
			t.Errorf("TimeToMonthContinuous(\"%s\"): want [%v] got [%v]",
				tt.rfc3339, tt.monthc, gotMonthc)
		}
	}
}

var monthFirstTests = []struct {
	year     int
	month    int
	want     string
	wantNext []string
	wantPrev []string
}{
	{2020, 7, "2020-07-01",
		[]string{"2020-08-01", "2020-09-01", "2020-10-01", "2020-11-01", "2020-12-01", "2021-01-01"},
		[]string{"2020-06-01", "2020-05-01", "2020-04-01", "2020-03-01", "2020-02-01", "2020-01-01", "2019-12-01"},
	},
}

func TestMonthFirst(t *testing.T) {
	for _, tt := range monthFirstTests {
		dt1 := time.Date(tt.year, time.Month(tt.month), 1, 0, 0, 0, 0, time.UTC)
		dt1Month := MonthBegin(dt1, 0)
		dt1MonthStr := dt1Month.Format(timeutil.RFC3339FullDate)
		if tt.want != dt1MonthStr {
			t.Errorf("MonthBegin(%v, %v): want [%v], got [%v]", dt1Month.Format(time.RFC3339),
				"0", tt.want, dt1MonthStr)
		}
		for i, want := range tt.wantNext {
			n := i + 1
			dtNext := MonthBegin(dt1, n)
			dtNextStr := dtNext.Format(timeutil.RFC3339FullDate)
			if want != dtNextStr {
				t.Errorf("MonthBegin(%v, %v): want [%v], got [%v]", dt1Month.Format(time.RFC3339),
					"0", want, dtNextStr)
			}
		}
		for i, want := range tt.wantPrev {
			n := (i + 1) * -1
			dtPrev := MonthBegin(dt1, n)
			dtPrevStr := dtPrev.Format(timeutil.RFC3339FullDate)
			if want != dtPrevStr {
				t.Errorf("MonthBegin(%v, %v): want [%v], got [%v]", dt1Month.Format(time.RFC3339),
					"0", want, dtPrevStr)
			}
		}
	}
}
