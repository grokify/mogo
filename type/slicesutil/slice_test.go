package slicesutil

import (
	"testing"
)

var splitMaxLengthTests = []struct {
	v  []uint
	s  uint
	n  uint
	l0 uint
}{
	{v: []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, s: 2, n: 5, l0: 9},
	{v: []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, s: 3, n: 4, l0: 10},
	{v: []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, s: 4, n: 3, l0: 9},
}

func TestSplitMaxLength(t *testing.T) {
	for _, tt := range splitMaxLengthTests {
		res := SplitMaxLength(tt.v, tt.s)
		if len(res) != int(tt.n) {
			t.Errorf("slicesutil.SplitMaxLength(%v, %d): mismatch on len: want (%d), got (%d)",
				tt.v, tt.s, tt.n, len(res))
		}
		if len(res) > 0 {
			l0 := len(res) - 1
			row := res[l0]
			if tt.l0 != row[0] {
				t.Errorf("slicesutil.SplitMaxLength(%v, %d): mismatch on value: want (%d), got (%d)",
					tt.v, tt.s, tt.l0, l0)
			}
		}
	}
}
