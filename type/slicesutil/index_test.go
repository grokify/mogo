package slicesutil

import (
	"testing"

	"github.com/grokify/mogo/math/mathutil"
)

var reverseIndexTests = []struct {
	n uint
	i uint
	r uint
}{
	{n: 5, i: 0, r: 4},
	{n: 5, i: 1, r: 3},
	{n: 5, i: 2, r: 2},
	{n: 5, i: 3, r: 1},
	{n: 5, i: 4, r: 0},
	{n: 5, i: 5, r: 4},
	{n: 5, i: 6, r: 3},
	{n: 5, i: 7, r: 2},
	{n: 5, i: 8, r: 1},
	{n: 5, i: 9, r: 0},

	{n: 4, i: 0, r: 3},
	{n: 4, i: 1, r: 2},
	{n: 4, i: 2, r: 1},
	{n: 4, i: 3, r: 0},
}

func TestReverseIndex(t *testing.T) {
	for _, tt := range reverseIndexTests {
		r := ReverseIndex(tt.n, tt.i)
		if r != tt.r {
			t.Errorf("slicesutil.ReverseIndex(%d, %d) (fwd) want (%d), got (%d)",
				tt.n, tt.i, tt.r, r)
		}
		if tt.i < tt.n {
			i := ReverseIndex(tt.n, r)
			if i != tt.i {
				t.Errorf("slicesutil.ReverseIndex(%d, %d) (rev) want (%d), got (%d)",
					tt.n, r, tt.i, i)
			}
		} else {
			i := ReverseIndex(tt.n, r)
			if i != mathutil.ModPyInt(tt.i, tt.n) {
				t.Errorf("slicesutil.ReverseIndex(%d, %d) (rev) want (%d), got (%d)",
					tt.n, r, tt.i, i)
			}
		}
	}
}
