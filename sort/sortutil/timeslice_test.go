package sortutil

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
