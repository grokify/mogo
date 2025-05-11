package month

import (
	"strings"
	"testing"

	"github.com/grokify/mogo/strconv/strconvutil"
)

var startEndDT6sTests = []struct {
	v   []int
	min int
	max int
}{
	{[]int{}, -1, -1},
	{[]int{202301}, 202301, 0},
	{[]int{202201, 202312}, 202201, 202312},
}

func TestStartEndDT6sTests(t *testing.T) {
	for _, tt := range startEndDT6sTests {
		gotSta, gotEnd := StartEndDT6s(tt.v)
		if gotSta != tt.min || gotEnd != tt.max {
			t.Errorf("month.StartEndDT6s(%s): want (%d, %d), got (%d, %d)",
				strings.Join(strconvutil.SliceItoaMore(tt.v, false, false), ","),
				tt.min, tt.max,
				gotSta, gotEnd)
		}
	}
}
