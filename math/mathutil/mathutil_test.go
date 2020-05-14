package mathutil

import (
	"strconv"
	"strings"
	"testing"
)

var divideInt64Tests = []struct {
	dividend  int64
	divisor   int64
	quotient  int64
	remainder int64
}{
	{15, 4, 3, 3},
}

func TestDivideInt64(t *testing.T) {
	for _, tt := range divideInt64Tests {
		quotient, remainder := DivideInt64(tt.dividend, tt.divisor)
		if tt.quotient != quotient || tt.remainder != remainder {
			t.Errorf("mathutil.DivideInt64(%d, %d) Mismatch: want [%d,%d], got [%d,%d]",
				tt.dividend, tt.divisor,
				tt.quotient, tt.remainder,
				quotient, remainder)
		}
	}
}

var minMaxTests = []struct {
	vals []int
	min  int
	max  int
}{
	{[]int{4, 7, 8}, 4, 8},
	{[]int{-4, -7, 8}, -7, 8},
}

func JoinIntsString(ints []int, sep string) string {
	strs := []string{}
	for _, val := range ints {
		strs = append(strs, strconv.Itoa(val))
	}
	return strings.Join(strs, sep)
}

func TestMinMax(t *testing.T) {
	for _, tt := range minMaxTests {
		strs := JoinIntsString(tt.vals, ",")
		// Test Int32
		int32s := IntsToInt32s(tt.vals)
		min32, max32 := MinMaxInt32(int32s...)
		if min32 != int32(tt.min) || max32 != int32(tt.max) {
			t.Errorf("mathutil.MinMaxInt32(%s) Mismatch: want [%d,%d], got [%d,%d]",
				strs, tt.min, tt.max, min32, max32)
		}
		if min32 >= 0 {
			// Test Uint
			uints := IntsToUints(tt.vals)
			minUint, maxUint := MinMaxUint(uints...)
			if minUint != uint(tt.min) || maxUint != uint(tt.max) {
				t.Errorf("mathutil.MinMaxUint(%s) Mismatch: want [%d,%d], got [%d,%d]",
					strs, tt.min, tt.max, minUint, maxUint)
			}
		}
	}
}
