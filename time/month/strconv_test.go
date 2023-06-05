package month

import (
	"strings"
	"testing"

	"github.com/grokify/mogo/type/slicesutil"
)

var startEndDT6sTests = []struct {
	v   []int32
	min int32
	max int32
}{
	{[]int32{}, -1, -1},
	{[]int32{202301}, 202301, 0},
	{[]int32{202201, 202312}, 202201, 202312},
}

func TestStartEndDT6sTests(t *testing.T) {
	for _, tt := range startEndDT6sTests {
		gotSta, gotEnd := StartEndDT6s(tt.v)
		if gotSta != tt.min || gotEnd != tt.max {
			t.Errorf("month.StartEndDT6s(%s): want (%d, %d), got (%d, %d)",
				strings.Join(slicesutil.Itoa(tt.v), ","),
				tt.min, tt.max,
				gotSta, gotEnd)
		}
	}
}
