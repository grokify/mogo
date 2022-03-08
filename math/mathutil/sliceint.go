package mathutil

import (
	"sort"
)

// SliceInt represets a slice of integers and provides functions on that slice.
type SliceInt struct {
	Elements []int
	Stats    SliceIntStats
}

// NewSliceInt creates and returns an empty SliceInt struct.
func NewSliceInt() SliceInt {
	sint := SliceInt{Elements: []int{}}
	return sint
}

// Append adds an element to the integer slice.
func (sint *SliceInt) Append(num int) {
	sint.Elements = append(sint.Elements, num)
}

// Len returns the number of items in the integer slice.
func (sint *SliceInt) Len() int {
	return len(sint.Elements)
}

// Sort sorts the elements in the integer slice.
func (sint *SliceInt) Sort() {
	sort.Ints(sint.Elements)
}

// Min returns the minimum element value in the integer slice.
func (sint *SliceInt) Min() (int, error) {
	if len(sint.Elements) == 0 {
		return 0, ErrEmptyList
	}
	if !sort.IntsAreSorted(sint.Elements) {
		sort.Ints(sint.Elements)
	}
	return sint.Elements[0], nil
}

// Max returns the maximum element value in the integer slice.
func (sint *SliceInt) Max() (int, error) {
	if len(sint.Elements) == 0 {
		return 0, ErrEmptyList
	}
	if !sort.IntsAreSorted(sint.Elements) {
		sort.Ints(sint.Elements)
	}
	return sint.Elements[len(sint.Elements)-1], nil
}

// Sum returns sum of all the elements in the integer slice.
func (sint *SliceInt) Sum() (int, error) {
	if len(sint.Elements) == 0 {
		return 0, ErrEmptyList
	}
	sum := int(0)
	for _, num := range sint.Elements {
		sum += num
	}
	return sum, nil
}

// Average is an alias for Mean.
func (sint *SliceInt) Average() (float64, error) {
	return sint.Mean()
}

// Mean returns the arithmetic mean of the integer slice.
func (sint *SliceInt) Mean() (float64, error) {
	if len(sint.Elements) == 0 {
		return 0, ErrEmptyList
	}
	sum, err := sint.Sum()
	if err != nil {
		return 0, err
	}
	return float64(sum) / float64(len(sint.Elements)), nil
}

// Median returns the median or middle value of the sorted integer slice.
func (sint *SliceInt) Median() (int, error) {
	if len(sint.Elements) == 0 {
		return 0, ErrEmptyList
	}
	if !sort.IntsAreSorted(sint.Elements) {
		sort.Ints(sint.Elements)
	}
	mid := int64(float64(len(sint.Elements)) / 2)
	return sint.Elements[mid], nil
}

// BuildStats builds a stats struct for current integer slice elements.
func (sint *SliceInt) BuildStats() (SliceIntStats, error) {
	stats := NewSliceIntStats()
	stats.Len = sint.Len()
	max, err := sint.Max()
	if err != nil {
		return stats, err
	}
	stats.Max = max
	min, err := sint.Min()
	if err != nil {
		return stats, err
	}
	stats.Min = min
	mean, err := sint.Mean()
	if err != nil {
		return stats, err
	}
	stats.Mean = mean
	median, err := sint.Median()
	if err != nil {
		return stats, err
	}
	stats.Median = median
	sum, err := sint.Sum()
	if err != nil {
		return stats, err
	}
	stats.Sum = sum
	sint.Stats = stats
	return stats, nil
}

// SliceIntStats represents a set of statistics for a set of integers.
type SliceIntStats struct {
	Len    int
	Max    int
	Mean   float64
	Median int
	Min    int
	Sum    int
}

// NewSliceIntStats returns a new initialized SliceIntStats struct.
func NewSliceIntStats() SliceIntStats {
	stats := SliceIntStats{
		Len:    0,
		Max:    0,
		Mean:   0,
		Median: 0,
		Min:    0,
		Sum:    0}
	return stats
}
