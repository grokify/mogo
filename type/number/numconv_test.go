package number

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/grokify/mogo/type/ordered"
)

var itou32Tests = []struct {
	v       int
	want    uint32
	wantErr bool
}{
	{v: -1, want: 0, wantErr: true},
	{v: 100, want: 100, wantErr: false},
	{v: int(^uint32(0)), want: math.MaxUint32, wantErr: false},
	{v: math.MaxInt, want: math.MaxUint32, wantErr: true},
	{v: math.MaxInt64, want: math.MaxUint32, wantErr: true},
	{v: 4294967295 + 1, want: math.MaxUint32, wantErr: true},
}

func TestItou32(t *testing.T) {
	for _, tt := range itou32Tests {
		got, err := Itou32(tt.v)
		if err != nil {
			if !tt.wantErr {
				t.Errorf("number.Uint32(%d) error (%s)",
					tt.v, err.Error())
			} else {
				continue
			}
		} else if tt.want != got {
			t.Errorf("number.Uint32(%d) want (%d) got (%d)",
				tt.v, tt.want, got)
		}
	}
}

var minMaxTests = []struct {
	vals []int
	min  int
	max  int
	minU uint
	maxU uint
}{
	{[]int{4, 7, 8}, 4, 8, 4, 8},
	{[]int{-4, -7, 8}, -7, 8, 0, 8},
}

func joinIntsString(ints []int, sep string) string {
	strs := []string{}
	for _, val := range ints {
		strs = append(strs, strconv.Itoa(val))
	}
	return strings.Join(strs, sep)
}

func TestMinMax(t *testing.T) {
	for _, tt := range minMaxTests {
		strs := joinIntsString(tt.vals, ",")
		// Test Int32
		int32s, err := Itoi32s(tt.vals)
		if err != nil {
			t.Errorf("mathutil.Itoi32s(%v) Error [%s]", tt.vals, err.Error())
		}
		min32, max32 := ordered.MinMax(int32s...)
		if int(min32) != tt.min || int(max32) != tt.max {
			t.Errorf("mathutil.MinMaxInt32(%s) Mismatch: want [%d,%d], got [%d,%d]",
				strs, tt.min, tt.max, min32, max32)
		}
		if min32 >= 0 {
			// Test Uint
			uints, err := Itous(tt.vals)
			if err != nil {
				t.Errorf("mathutil.Itous(%v) Error [%s]", tt.vals, err.Error())
			}
			minUint, maxUint := ordered.MinMax(uints...)
			if minUint != tt.minU || maxUint != tt.maxU {
				t.Errorf("mathutil.MinMaxUint(%s) Mismatch: want [%d,%d], got [%d,%d]",
					strs, tt.min, tt.max, minUint, maxUint)
			}
		}
	}
}
