package mathutil

import (
	"testing"
)

var intLenTests = []struct {
	val int
	len int
}{
	{-1000, 5},
	{-1, 2},
	{0, 1},
	{1, 1},
	{1000, 4},
}

func TestIntLen(t *testing.T) {
	for _, tt := range intLenTests {
		intLen := IntLen(tt.val)
		if intLen != tt.len {
			t.Errorf("mathutil.IntLen(%d) Mismatch: want [%d], got [%d]",
				tt.val, tt.len, intLen)
		}
	}
}

func TestMaxInt63(t *testing.T) {
	want := 4611686018427387903
	if MaxInt63 != want {
		t.Errorf("mathutil.MaxInt63 Mismatch: want [%d], got [%d]",
			want, MaxInt63)
	}
}
