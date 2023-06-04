package mathutil

import (
	"testing"
)

var divideInt64Tests = []struct {
	dividend  int64
	divisor   int64
	quotient  int64
	remainder int64
}{
	{15, 4, 3, 3},
	{8, 1, 8, 0},
	{8, 2, 4, 0},
	{8, 3, 2, 2},
	{8, 4, 2, 0},
	{8, 5, 1, 3},
	{8, 6, 1, 2},
	{8, 7, 1, 1},
	{8, 8, 1, 0},
	{-5, 3, -999, 1},
	{-4, 3, -999, 2},
	{-3, 3, -999, 0},
	{-2, 3, -999, 1},
	{-1, 3, -999, 2},
	{-0, 3, -999, 0},
	{1, 3, -999, 1},
	{2, 3, -999, 2},
	{3, 3, -999, 0},
	{4, 3, -999, 1},
}

func TestDivideInt64(t *testing.T) {
	for _, tt := range divideInt64Tests {
		if tt.quotient > 0 {
			quotient, remainder := DivideInt64(tt.dividend, tt.divisor)
			if tt.quotient != quotient || tt.remainder != remainder {
				t.Errorf("mathutil.DivideInt64(%d, %d) Mismatch: want [%d,%d], got [%d,%d]",
					tt.dividend, tt.divisor,
					tt.quotient, tt.remainder,
					quotient, remainder)
			}
		}
		/*
			modInt := ModInt(int(tt.dividend), int(tt.divisor))
			if modInt != int(tt.remainder) {
				t.Errorf("mathutil.ModInt(%d, %d) Mismatch: want [%d], got [%d]",
					tt.dividend, tt.divisor,
					tt.remainder, modInt)
			}
			modInt64 := ModInt64(tt.dividend, tt.divisor)
			if modInt64 != tt.remainder {
				t.Errorf("mathutil.ModInt64(%d, %d) Mismatch: want [%d], got [%d]",
					tt.dividend, tt.divisor,
					tt.remainder, modInt64)
			}
		*/
		modPy := ModPyInt(tt.dividend, tt.divisor)
		if modPy != tt.remainder {
			t.Errorf("mathutil.PyMod(%d, %d) Mismatch: want [%d], got [%d]",
				tt.dividend, tt.divisor,
				tt.remainder, modPy)
		}
	}
}

var modTests = []struct {
	input int
	mod   int
}{
	{-8, 0},
	{-7, 1},
	{-6, 2},
	{-5, 3},
	{-4, 0},
	{-3, 1},
	{-2, 2},
	{-1, 3},
	{0, 0},
	{1, 1},
	{2, 2},
	{3, 3},
	{4, 0},
	{5, 1},
	{6, 2},
	{7, 3},
	{8, 0},
	{9, 1},
	{10, 2},
	{11, 3},
}

func TestMod(t *testing.T) {
	for _, tt := range modTests {
		try := ModInt(tt.mod, int(4))
		if try != tt.mod {
			t.Errorf("mathutil.ModInt(%d, %d) Mismatch: want [%d], got [%d]",
				tt.mod, 4, tt.mod, try)
		}
	}
}
