package bigutil

import (
	"math/big"
	"testing"
)

var powBigIntTests = []struct {
	x    int64
	y    int64
	want string
}{
	{2, 0, "1"},
	{2, 1, "2"},
	{2, 2, "4"},
	{2, 3, "8"},
	{2, 4, "16"}}

func TestPowInt(t *testing.T) {
	for _, tt := range powBigIntTests {
		x := big.NewInt(tt.x)
		y := big.NewInt(tt.y)
		pow := PowInt(x, y)

		if pow.String() != tt.want {
			t.Errorf("bigutil.PowInt(%v,%v): want [%v], got [%v]", tt.x, tt.y, tt.want, pow.String())
		}
	}
}
