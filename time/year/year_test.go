package year

import (
	"testing"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

var timeSeriesYearTests = []struct {
	input  []string
	series []string
}{
	{
		[]string{
			"2010-01-01T00:00:00Z",
			"2000-01-01T00:00:00Z"},
		[]string{
			"2000-01-01T00:00:00Z",
			"2001-01-01T00:00:00Z",
			"2002-01-01T00:00:00Z",
			"2003-01-01T00:00:00Z",
			"2004-01-01T00:00:00Z",
			"2005-01-01T00:00:00Z",
			"2006-01-01T00:00:00Z",
			"2007-01-01T00:00:00Z",
			"2008-01-01T00:00:00Z",
			"2009-01-01T00:00:00Z",
			"2010-01-01T00:00:00Z"},
	}}

func TestTimeSeriesYear(t *testing.T) {
	for _, tt := range timeSeriesYearTests {
		input, err := timeutil.ParseTimes(time.RFC3339, tt.input)
		if err != nil {
			t.Errorf("year.TestTimeSeriesYear cannot parse [%v] Error: [%s]", tt.input, err.Error())
		}
		seriesWant, err := timeutil.ParseTimes(time.RFC3339, tt.series)
		if err != nil {
			t.Errorf("year.TestTimeSeriesYear cannot parse [%v] Error: [%s]", tt.series, err.Error())
		}
		seriesTry := TimeSeriesYear(true, input...)
		if !seriesTry.Equal(seriesWant) {
			t.Errorf("year.TestTimeSeriesYear series not equal: want [%v] try [%v]", seriesWant, seriesTry)
		}
	}
}
