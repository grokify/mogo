package mathutil

import (
	"errors"
	"fmt"

	pkgerr "github.com/pkg/errors"
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
func (rf *RangeFloat64) CellIndexForValue(v float64) (int32, error) {
	err := rf.isInitialized()
	if err != nil {
		return int32(0), err
	}
	if v < rf.Min || v > rf.Max {
		return int32(0), errors.New(fmt.Sprintf("Value (%v) out of range [%v,%v]", v, rf.Min, rf.Max))
	}
	rf.iter = 0
	return rf.binarySearch(v, 0, rf.Cells-1)
}

func (rf *RangeFloat64) binarySearch(v float64, l, r int32) (int32, error) {
	rf.iter += 1
	if rf.iter > MaxTries {
		return int32(0), errors.New(fmt.Sprintf("Too many (%v) binary search iterations.", MaxTries))
	}

	m := int32(float32(l) + (float32(r)-float32(l))/2.0)
	min, max, err := rf.CellMinMax(m)
	if err != nil {
		return int32(0), pkgerr.Wrap(err, "CellMinMax() failed")
	}
	//fmt.Printf("{\"iter\":%v,\"val\":%v,\"l\":%v,\"r\":%v,\"m\":%v,\"minv\":%v,\"maxv\":%v}\n", rf.iter, v, l, r, m, min, max)

	if v >= min && (v < max || max == rf.Max) {
		return m, nil
	}
	if v < min {
		return rf.binarySearch(v, l, m-1)
	}
	return rf.binarySearch(v, m+1, r)
}

func (rf *RangeFloat64) isInitialized() error {
	if rf.Min > rf.Max {
		return errors.New(fmt.Sprintf("Start (%v) is less than End (%v)", rf.Min, rf.Max))
	} else if rf.Cells <= 0 {
		return errors.New(fmt.Sprintf("Num cells is <= 0 (%v)", rf.Cells))
	}
	return nil
}

func (rf *RangeFloat64) cellRange() (float64, error) {
	err := rf.isInitialized()
	if err != nil {
		return float64(0), err
	}
	return ((rf.Max - rf.Min) / float64(rf.Cells)), nil
}

func (rf *RangeFloat64) CellMinMax(idx int32) (float64, float64, error) {
	err := rf.isInitialized()
	if err != nil {
		return float64(0), float64(0), err
	}
	cellRange, err := rf.cellRange()
	if err != nil {
		return float64(0), float64(0), err
	}
	cellMin := float64(idx) * cellRange
	cellMax := rf.Max
	if idx < (rf.Cells - 1) {
		cellMax = (float64(idx) + 1) * cellRange
	}
	return cellMin, cellMax, nil
}
