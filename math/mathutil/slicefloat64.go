package mathutil

import (
	"errors"
	"sort"
)

var ErrEmptyList = errors.New("list is empty")

// SliceFloat64 represets a slice of integers and provides functions on that slice.
type SliceFloat64 struct {
	Elements []float64
	Stats    SliceFloat64Stats
}

// NewSliceFloat64 creates and returns an empty SliceFloat64 struct.
func NewSliceFloat64() SliceFloat64 {
	sf64 := SliceFloat64{Elements: []float64{}}
	return sf64
}

// Append adds an element to the float64 slice.
func (sf64 *SliceFloat64) Append(num float64) {
	sf64.Elements = append(sf64.Elements, num)
}

// Len returns the number of items in the float64 slice.
func (sf64 *SliceFloat64) Len() int {
	return len(sf64.Elements)
}

// Sort sorts the elements in the float64 slice.
func (sf64 *SliceFloat64) Sort() {
	sort.Float64s(sf64.Elements)
}

// Min returns the minimum element value in the float64 slice.
func (sf64 *SliceFloat64) Min() (float64, error) {
	if len(sf64.Elements) == 0 {
		return 0, ErrEmptyList
	}
	if !sort.Float64sAreSorted(sf64.Elements) {
		sort.Float64s(sf64.Elements)
	}
	return sf64.Elements[0], nil
}

// Max returns the maximum element value in the float64 slice.
func (sf64 *SliceFloat64) Max() (float64, error) {
	if len(sf64.Elements) == 0 {
		return 0, ErrEmptyList
	}
	if !sort.Float64sAreSorted(sf64.Elements) {
		sort.Float64s(sf64.Elements)
	}
	return sf64.Elements[len(sf64.Elements)-1], nil
}

// Sum returns sum of all the elements in the float64 slice.
func (sf64 *SliceFloat64) Sum() (float64, error) {
	if len(sf64.Elements) == 0 {
		return 0, ErrEmptyList
	}
	sum := float64(0)
	for _, num := range sf64.Elements {
		sum += num
	}
	return sum, nil
}

// Average is an alias for Mean.
func (sf64 *SliceFloat64) Average() (float64, error) {
	return sf64.Mean()
}

// Mean returns the arithmetic mean of the float64 slice.
func (sf64 *SliceFloat64) Mean() (float64, error) {
	if len(sf64.Elements) == 0 {
		return 0, ErrEmptyList
	}
	sum, err := sf64.Sum()
	if err != nil {
		return 0, err
	}
	return sum / float64(len(sf64.Elements)), nil
}

// Median returns the median or middle value of the sorted float64 slice.
func (sf64 *SliceFloat64) Median() (float64, error) {
	if len(sf64.Elements) == 0 {
		return 0, ErrEmptyList
	}
	if !sort.Float64sAreSorted(sf64.Elements) {
		sort.Float64s(sf64.Elements)
	}
	mid := int64(float64(len(sf64.Elements)) / 2)
	return sf64.Elements[mid], nil
}

// BuildStats builds a stats struct for current float64 slice elements.
func (sf64 *SliceFloat64) BuildStats() (SliceFloat64Stats, error) {
	stats := NewSliceFloat64Stats()
	stats.Len = sf64.Len()
	max, err := sf64.Max()
	if err != nil {
		return stats, err
	}
	stats.Max = max
	min, err := sf64.Min()
	if err != nil {
		return stats, err
	}
	stats.Min = min
	mean, err := sf64.Mean()
	if err != nil {
		return stats, err
	}
	stats.Mean = mean
	median, err := sf64.Median()
	if err != nil {
		return stats, err
	}
	stats.Median = median
	sum, err := sf64.Sum()
	if err != nil {
		return stats, err
	}
	stats.Sum = sum
	sf64.Stats = stats
	return stats, nil
}

// SliceFloat64Stats represents a set of statistics for a set of float64s.
type SliceFloat64Stats struct {
	Len    int
	Max    float64
	Mean   float64
	Median float64
	Min    float64
	Sum    float64
}

// NewSliceFloat64Stats returns a new initialized SliceFloat64Stats struct.
func NewSliceFloat64Stats() SliceFloat64Stats {
	stats := SliceFloat64Stats{
		Len:    0,
		Max:    0,
		Mean:   0,
		Median: 0,
		Min:    0,
		Sum:    0}
	return stats
}
