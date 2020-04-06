package mathutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	pkgerr "github.com/pkg/errors"
)

// RangeFloat64 creates a range with a fixed set of cells and
// returns the cell for a certain value.
type RangeInt64 struct {
	Min   int64
	Max   int64
	Cells int32
	iter  int32
}

// CellIndexForValue returns a cell index for a requested value.
func (rng *RangeInt64) CellIndexForValue(v int64) (int32, error) {
	err := rng.isInitialized()
	if err != nil {
		return int32(0), err
	}
	if v < rng.Min || v > rng.Max {
		return int32(0), fmt.Errorf("RangeInt64 Value (%v) out of range [%v,%v]", v, rng.Min, rng.Max)
	}
	rng.iter = 0
	return rng.binarySearch(v, 0, rng.Cells-1)
}

func (rng *RangeInt64) binarySearch(v int64, l, r int32) (int32, error) {
	rng.iter += 1
	if rng.iter > MaxTries {
		return int32(0), fmt.Errorf("RangeInt64 Too many (%v) binary search iterations.", MaxTries)
	}

	m := int32(float32(l) + (float32(r)-float32(l))/2.0)
	min, max, err := rng.CellMinMax(m)
	if err != nil {
		return int32(0), pkgerr.Wrap(err, "CellMinMax() failed")
	}
	//fmt.Printf("{\"iter\":%v,\"val\":%v,\"l\":%v,\"r\":%v,\"m\":%v,\"minv\":%v,\"maxv\":%v}\n", rng.iter, v, l, r, m, min, max)

	if v >= min && (v < max || max == rng.Max) {
		return m, nil
	}
	if v < min {
		return rng.binarySearch(v, l, m-1)
	}
	return rng.binarySearch(v, m+1, r)
}

func (rng *RangeInt64) isInitialized() error {
	if rng.Min > rng.Max {
		return fmt.Errorf("Start (%v) is less than End (%v)", rng.Min, rng.Max)
	} else if rng.Cells <= 0 {
		return fmt.Errorf("Num cells is <= 0 (%v)", rng.Cells)
	}
	return nil
}

func (rng *RangeInt64) CellRange() (int64, error) {
	err := rng.isInitialized()
	if err != nil {
		return int64(0), err
	}
	return int64(float64(rng.Max-rng.Min) / float64(rng.Cells)), nil
}

func (rng *RangeInt64) CellMinMax(idx int32) (int64, int64, error) {
	err := rng.isInitialized()
	if err != nil {
		return int64(0), int64(0), err
	}
	cellRange, err := rng.CellRange()
	if err != nil {
		return int64(0), int64(0), err
	}
	cellMin := int64(idx)*cellRange + rng.Min
	cellMax := rng.Max
	if idx < (rng.Cells - 1) {
		cellMax = (int64(idx)+1)*cellRange + rng.Min
	}
	return cellMin, cellMax, nil
}

// DivmodInt64 from https://stackoverflow.com/a/43945812
func DivmodInt64(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

// PrettyTicks returns a slice of integers that start
// lower and end higher than the supplied range. This
// is intended to be used for chart axis.
func PrettyTicks(estimatedTickCount int, low, high int64) []int64 {
	ticks := []int64{}
	if low > high {
		tmp := low
		low = high
		high = tmp
	}
	diffRaw := high - low
	tickSize := float64(diffRaw) / float64(estimatedTickCount)
	tickSizedRounded := FloorMostSignificant(int64(tickSize))
	lowFloor := FloorMostSignificant(int64(low))

	ticks = append(ticks, lowFloor)
	for ticks[len(ticks)-1] < high {
		ticks = append(ticks, ticks[len(ticks)-1]+tickSizedRounded)
	}
	return ticks
}

// FloorMostSignificant returns number with a single
// significant digit followed by zeros. The value returned
// is always lower than the supplied value.
func FloorMostSignificant(valOriginal int64) int64 {
	// see here for additional discussion
	// https://stackoverflow.com/questions/202302/
	if valOriginal == 0 {
		return 0
	}
	valPositive := valOriginal
	isNegative := false
	if valOriginal < 0 {
		valPositive = -1 * valOriginal
		isNegative = true
	}
	valStr := fmt.Sprintf("%d", valPositive)
	valLen := len(fmt.Sprintf("%d", valPositive))
	var final int64

	// Math power approach
	mostSig, err := strconv.Atoi(valStr[0:1])
	if err != nil {
		panic(errors.Wrap(err, "mathutil.FloorMostSignificant"))
	}
	if isNegative {
		final = -1 * int64(mostSig+1) * int64(math.Pow10(valLen-1))
	} else {
		final = int64(mostSig) * int64(math.Pow10(valLen-1))
	}
	// String approach
	if 1 == 0 {
		vals := make([]string, valLen)
		for i := 0; i < valLen; i++ {
			if i == 0 {
				vals[i] = valStr[0:1]
			} else {
				vals[i] = "0"
			}
		}
		intStr := strings.Join(vals, "")
		num, err := strconv.Atoi(intStr)
		if err != nil {
			panic(errors.Wrap(err, "mathutil.FloorMostSignificant"))
		}
		final = int64(num)
	}
	return final
}
