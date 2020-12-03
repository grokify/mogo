package bigint

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
		pow := Pow(x, y)

		if pow.String() != tt.want {
			t.Errorf("bigutil.PowInt(%v,%v): want [%v], got [%v]", tt.x, tt.y, tt.want, pow.String())
		}
	}
}

var intToBaseXStringTests = []struct {
	base int
	dec  int
	val  string
}{
	{16, 15, "f"},
	{32, 0, "0"},
	{32, 1, "1"},
	{32, 12, "c"}}

const base36Dictionary = "0123456789abcdefghijklmnopqrstuvwxyz"

func TestIntToBaseXString(t *testing.T) {
	for _, tt := range intToBaseXStringTests {
		try := IntToBaseXString(tt.base, tt.dec)

		if try != tt.val {
			t.Errorf("bigint.IntToBaseXString(%v,%v): want [%v], got [%v]",
				tt.base, tt.dec, tt.val, try)
		}
		try2 := string([]rune(base36Dictionary)[tt.dec])
		if try2 != tt.val {
			t.Errorf("bigint.IntToBaseXString(%v,%v): want [%s], got [%s]",
				tt.base, tt.dec, tt.val, try2)
		}
	}
}
