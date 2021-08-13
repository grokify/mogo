package timeslice

import (
	"sort"
	"strings"
	"testing"
	"time"
)

var timeSliceTests = []struct {
	rfc3339Times []string
	wantSorted   []string
	wantDeduped  []string
}{
	{
		[]string{
			"2021-01-01T01:01:01Z", "2021-01-01T01:01:01Z",
			"2021-01-01T00:00:00Z", "2021-01-01T02:02:02Z"},
		[]string{
			"2021-01-01T00:00:00Z", "2021-01-01T01:01:01Z",
			"2021-01-01T01:01:01Z", "2021-01-01T02:02:02Z"},
		[]string{"2021-01-01T01:01:01Z",
			"2021-01-01T00:00:00Z", "2021-01-01T02:02:02Z"}},
}

func TestTimeSlice(t *testing.T) {
	for _, tt := range timeSliceTests {
		timeSlice, err := ParseTimeSlice(time.RFC3339, tt.rfc3339Times)
		if err != nil {
			t.Errorf("time slice did not parse as RFC-3339 [%s] error [%s]",
				strings.Join(tt.rfc3339Times, ","), err.Error())
		}
		wantDedupedTimeSlice, err := ParseTimeSlice(time.RFC3339, tt.wantDeduped)
		if err != nil {
			t.Errorf("time slice did not parse as RFC-3339 [%s] error [%s]",
				strings.Join(tt.wantDeduped, ","), err.Error())
		}
		dedupedTimeSlice := timeSlice.Dedupe()
		if !wantDedupedTimeSlice.Equal(dedupedTimeSlice) {
			t.Errorf("timeSlice.Dedupe FAIL want [%s] got [%s]",
				strings.Join(tt.wantDeduped, ","),
				strings.Join(dedupedTimeSlice.Format(time.RFC3339), ","))
		}
		wantSortedTimeSlice, err := ParseTimeSlice(time.RFC3339, tt.wantSorted)
		if err != nil {
			t.Errorf("time slice did not parse as RFC-3339 [%s] error [%s]",
				strings.Join(tt.wantSorted, ","), err.Error())
		}
		sort.Sort(timeSlice)
		if !wantSortedTimeSlice.Equal(timeSlice) {
			t.Errorf("timeSlice.Dedupe FAIL want [%s] got [%s]",
				strings.Join(tt.wantDeduped, ","),
				strings.Join(timeSlice.Format(time.RFC3339), ","))
		}
	}
}

var timeRangeTests = []struct {
	timeRange []string
	timeTest  string
	inclusive bool
	testLower bool
	testUpper bool
	wantLower string
	wantUpper string
}{
	{
		timeRange: []string{
			"2021-06-01T01:01:01Z", "2020-06-01T01:01:01Z",
			"2021-12-01T01:01:01Z"},
		timeTest:  "2021-07-01T01:01:01Z",
		inclusive: true,
		testLower: true,
		testUpper: true,
		wantLower: "2021-06-01T01:01:01Z",
		wantUpper: "2021-12-01T01:01:01Z",
	},
	{
		timeRange: []string{
			"2021-06-01T01:01:01Z", "2020-06-01T01:01:01Z",
			"2021-12-01T01:01:01Z"},
		timeTest:  "2021-06-01T01:01:01Z",
		inclusive: true,
		testLower: true,
		testUpper: true,
		wantLower: "2021-06-01T01:01:01Z",
		wantUpper: "2021-06-01T01:01:01Z",
	},
	{
		timeRange: []string{
			"2021-06-01T01:01:01Z", "2020-06-01T01:01:01Z",
			"2021-12-01T01:01:01Z"},
		timeTest:  "2020-07-01T01:01:01Z",
		inclusive: true,
		testLower: true,
		testUpper: true,
		wantLower: "2020-06-01T01:01:01Z",
		wantUpper: "2021-06-01T01:01:01Z",
	},
	{
		timeRange: []string{
			"2021-06-01T01:01:01Z", "2020-06-01T01:01:01Z",
			"2021-12-01T01:01:01Z"},
		timeTest:  "2019-07-01T01:01:01Z",
		inclusive: true,
		testLower: false,
		testUpper: true,
		wantUpper: "2020-06-01T01:01:01Z",
	},
}

func TestTimeRange(t *testing.T) {
	for i, tt := range timeRangeTests {
		timeSlice, err := ParseTimeSlice(time.RFC3339, tt.timeRange)
		if err != nil {
			t.Errorf("time slice did not parse as RFC-3339 [%s] error [%s]",
				strings.Join(tt.timeRange, ","), err.Error())
		}
		timeTest, err := time.Parse(time.RFC3339, tt.timeTest)
		if err != nil {
			t.Errorf("time did not parse as RFC-3339 [%s] error [%s]",
				tt.timeTest, err.Error())
		}
		if tt.testLower {
			wantLower, err := time.Parse(time.RFC3339, tt.wantLower)
			if err != nil {
				t.Errorf("time did not parse as RFC-3339 [%s] error [%s]",
					tt.wantLower, err.Error())
			}
			tryLower, err := timeSlice.RangeLower(timeTest, tt.inclusive)
			if err != nil {
				t.Errorf("timeSlice.RangeLower try [%s] error [%s]",
					tt.timeTest, err.Error())
			}
			if !tryLower.Equal(wantLower) {
				t.Errorf("timeSlice.RangeLower wrong value: using [%s] got [%s] want [%s] range [%s]",
					timeTest.Format(time.RFC3339),
					tryLower.Format(time.RFC3339),
					wantLower.Format(time.RFC3339),
					strings.Join(timeSlice.Format(time.RFC3339), ","))
			}
		}
		if tt.testUpper {
			wantUpper, err := time.Parse(time.RFC3339, tt.wantUpper)
			if err != nil {
				t.Errorf("time did not parse as RFC-3339 [%s] error [%s]",
					tt.wantUpper, err.Error())
			}
			tryUpper, err := timeSlice.RangeUpper(timeTest, tt.inclusive)
			if err != nil {
				t.Errorf("timeSlice.RangeUpper try [%s] error [%s]",
					tt.timeTest, err.Error())
			}
			if !tryUpper.Equal(wantUpper) {
				t.Errorf("timeSlice.RangeUpper wrong value: using [%s] got [%s] want [%s] range [%s]",
					timeTest.Format(time.RFC3339),
					tryUpper.Format(time.RFC3339),
					wantUpper.Format(time.RFC3339),
					strings.Join(timeSlice.Format(time.RFC3339), ","))
			}
		}
	}
}
