package mathutil

import (
	"fmt"

	"github.com/grokify/mogo/errors/errorsutil"
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
		return int32(0), fmt.Errorf("struct RangeInt64 value (%v) out of range [%v,%v]", v, rng.Min, rng.Max)
	}
	rng.iter = 0
	return rng.binarySearch(v, 0, rng.Cells-1)
}

func (rng *RangeInt64) binarySearch(v int64, l, r int32) (int32, error) {
	rng.iter += 1
	if rng.iter > MaxTries {
		return int32(0), fmt.Errorf("struct RangeInt64 too many (%v) binary search iterations", MaxTries)
	}

	m := int32(float32(l) + (float32(r)-float32(l))/2.0)
	min, max, err := rng.CellMinMax(m)
	if err != nil {
		return int32(0), errorsutil.Wrap(err, "func CellMinMax() failed")
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
		return fmt.Errorf("start (%v) is less than End (%v)", rng.Min, rng.Max)
	} else if rng.Cells <= 0 {
		return fmt.Errorf("num cells is <= 0 (%v)", rng.Cells)
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
