package mathutil

import (
	"fmt"

	"github.com/grokify/mogo/errors/errorsutil"
)

var MaxTries = int32(100000)

// RangeFloat64 creates a range with a fixed set of cells and
// returns the cell for a certain value.
type RangeFloat64 struct {
	Min   float64
	Max   float64
	Cells int32
	iter  int32
}

// CellIndexForValue returns a cell index for a requested value.
func (rng *RangeFloat64) CellIndexForValue(v float64) (int32, error) {
	err := rng.isInitialized()
	if err != nil {
		return int32(0), err
	}
	if v < rng.Min || v > rng.Max {
		return int32(0), fmt.Errorf("Value (%v) out of range [%v,%v]", v, rng.Min, rng.Max)
	}
	rng.iter = 0
	return rng.binarySearch(v, 0, rng.Cells-1)
}

func (rng *RangeFloat64) binarySearch(v float64, l, r int32) (int32, error) {
	rng.iter += 1
	if rng.iter > MaxTries {
		return int32(0), fmt.Errorf("Too many (%v) binary search iterations.", MaxTries)
	}

	m := int32(float32(l) + (float32(r)-float32(l))/2.0)
	min, max, err := rng.CellMinMax(m)
	if err != nil {
		return int32(0), errorsutil.Wrap(err, "CellMinMax() failed")
	}
	//fmt.Printf("{\"iter\":%v,\"val\":%v,\"l\":%v,\"r\":%v,\"m\":%v,\"minv\":%v,\"maxv\":%v}\n", rf.iter, v, l, r, m, min, max)

	if v >= min && (v < max || max == rng.Max) {
		return m, nil
	}
	if v < min {
		return rng.binarySearch(v, l, m-1)
	}
	return rng.binarySearch(v, m+1, r)
}

func (rng *RangeFloat64) isInitialized() error {
	if rng.Min > rng.Max {
		return fmt.Errorf("Start (%v) is less than End (%v)", rng.Min, rng.Max)
	} else if rng.Cells <= 0 {
		return fmt.Errorf("Num cells is <= 0 (%v)", rng.Cells)
	}
	return nil
}

func (rng *RangeFloat64) cellRange() (float64, error) {
	err := rng.isInitialized()
	if err != nil {
		return float64(0), err
	}
	return ((rng.Max - rng.Min) / float64(rng.Cells)), nil
}

func (rng *RangeFloat64) CellMinMax(idx int32) (float64, float64, error) {
	err := rng.isInitialized()
	if err != nil {
		return float64(0), float64(0), err
	}
	cellRange, err := rng.cellRange()
	if err != nil {
		return float64(0), float64(0), err
	}
	cellMin := float64(idx)*cellRange + rng.Min
	cellMax := rng.Max
	if idx < (rng.Cells - 1) {
		cellMax = (float64(idx)+1)*cellRange + rng.Min
	}
	return cellMin, cellMax, nil
}
