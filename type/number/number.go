package number

import (
	"errors"
	"math"
	"sort"
	"strconv"

	"golang.org/x/exp/constraints"
)

// SliceToFloat64 converts a slice of integers or floats to float64.
func SliceToFloat64[T constraints.Integer | constraints.Float](src []T) []float64 {
	var out []float64
	for _, in := range src {
		out = append(out, float64(in))
	}
	return out
}

func IntSliceDedupe(elems []int, sortElems bool) []int {
	if len(elems) == 0 {
		return elems
	}
	var deduped []int
	mapInts := map[int]int{}
	for _, el := range elems {
		if _, ok := mapInts[el]; !ok {
			deduped = append(deduped, el)
			mapInts[el] = 1
		}
	}
	if sortElems {
		sort.Ints(deduped)
	}
	return deduped
}

func IntLength[I int | int8 | int16 | int32 | int64](num I) uint {
	return uint(len(strconv.Itoa(int(num))))
}

var ErrOverflow = errors.New("integer overflow")

func Int32(i int) (int32, error) {
	if i > math.MaxInt32 || i < math.MinInt32 {
		return 0, ErrOverflow
	}
	return int32(i), nil
}

func Int16(i int) (int16, error) {
	if i > math.MaxInt16 || i < math.MinInt16 {
		return 0, ErrOverflow
	}
	return int16(i), nil
}

func Int8(i int) (int8, error) {
	if i > math.MaxInt8 || i < math.MinInt8 {
		return 0, ErrOverflow
	}
	return int8(i), nil
}

/*
// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}
*/

// SliceIntBuildBeginEnd returns of slice of integers from beginning to end, inclusive.
func SliceIntBuildBeginEnd(beg, end int) []int {
	if beg > end {
		return []int{}
	} else if beg == end {
		return []int{beg}
	}
	var out []int
	for i := beg; i <= end; i++ {
		out = append(out, i)
	}
	return out
}
