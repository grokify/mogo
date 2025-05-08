package month

import (
	"testing"
	"time"

	"github.com/grokify/mogo/time/timeutil"
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
	year  uint32
	month uint32
	want  string
}{
	{uint32(1), uint32(1), "002T"},
	{uint32(500), uint32(3), "12KZ"},
	{uint32(1000), uint32(6), "255Y"},
	{uint32(2010), uint32(8), "4B3K"},
	{uint32(2019), uint32(8), "4BSK"},
	{uint32(2038), uint32(1), "4D95"},
	{uint32(9999), uint32(12), "LFJC"},
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
	year    uint32
	month   uint32
	monthc  uint32
	rfc3339 string
}{
	{uint32(0), uint32(1), uint32(1), "0000-01-01T00:00:00Z"},
	{uint32(0), uint32(12), uint32(12), "0000-12-01T00:00:00Z"},
	{uint32(1), uint32(1), uint32(13), "0001-01-01T00:00:00Z"},
	{uint32(1), uint32(12), uint32(24), "0001-12-01T00:00:00Z"},
	{uint32(2), uint32(1), uint32(25), "0002-01-01T00:00:00Z"},
	{uint32(3), uint32(12), uint32(48), "0003-12-01T00:00:00Z"},
	{uint32(4), uint32(1), uint32(49), "0004-01-01T00:00:00Z"},
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
		gotMonthc, err := TimeToMonthContinuous(wantDt)
		if err != nil {
			t.Errorf("TimeToMonthContinuous(\"%v\") error [%v]",
				wantDt.Format(time.RFC3339), err.Error())
		}
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
		dt1Month := MonthStart(dt1, 0)
		dt1MonthStr := dt1Month.Format(timeutil.RFC3339FullDate)
		if tt.want != dt1MonthStr {
			t.Errorf("MonthStart(%v, %v): want [%v], got [%v]", dt1Month.Format(time.RFC3339),
				"0", tt.want, dt1MonthStr)
		}
		for i, want := range tt.wantNext {
			n := i + 1
			dtNext := MonthStart(dt1, n)
			dtNextStr := dtNext.Format(timeutil.RFC3339FullDate)
			if want != dtNextStr {
				t.Errorf("MonthStart(%v, %v): want [%v], got [%v]", dt1Month.Format(time.RFC3339),
					"0", want, dtNextStr)
			}
		}
		for i, want := range tt.wantPrev {
			n := (i + 1) * -1
			dtPrev := MonthStart(dt1, n)
			dtPrevStr := dtPrev.Format(timeutil.RFC3339FullDate)
			if want != dtPrevStr {
				t.Errorf("MonthStart(%v, %v): want [%v], got [%v]", dt1Month.Format(time.RFC3339),
					"0", want, dtPrevStr)
			}
		}
	}
}

var timesStartsMonthTests = []struct {
	input  []string
	series []string
}{
	{
		[]string{
			"2000-11-01T00:00:00Z",
			"2001-04-01T00:00:00Z",
		},
		[]string{
			"2000-11-01T00:00:00Z",
			"2000-12-01T00:00:00Z",
			"2001-01-01T00:00:00Z",
			"2001-02-01T00:00:00Z",
			"2001-03-01T00:00:00Z",
			"2001-04-01T00:00:00Z",
		},
	},
	{
		[]string{
			"2000-11-15T00:00:00Z",
			"2001-04-20T00:00:00Z",
		},
		[]string{
			"2000-11-01T00:00:00Z",
			"2000-12-01T00:00:00Z",
			"2001-01-01T00:00:00Z",
			"2001-02-01T00:00:00Z",
			"2001-03-01T00:00:00Z",
			"2001-04-01T00:00:00Z",
		},
	},
}

func TestTimesStartsMonth(t *testing.T) {
	for _, tt := range timesStartsMonthTests {
		input, err := timeutil.ParseTimes(time.RFC3339, tt.input)
		if err != nil {
			t.Errorf("year.TestTimeSeriesMonth cannot parse [%v] Error: [%s]", tt.input, err.Error())
		}
		seriesWant, err := timeutil.ParseTimes(time.RFC3339, tt.series)
		if err != nil {
			t.Errorf("year.TestTimeSeriesMonth cannot parse [%v] Error: [%s]", tt.series, err.Error())
		}
		seriesTry := TimesMonthStarts(input...)
		if !seriesTry.Equal(seriesWant) {
			t.Errorf("year.TimeSeriesMonth series not equal: want [%v] try [%v]", seriesWant, seriesTry)
		}
	}
}
