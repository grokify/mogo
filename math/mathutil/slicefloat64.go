package mathutil

import (
	"errors"
	"sort"
)

type SliceFloat64 struct {
	Elements []float64
}

func NewSliceFloat64() SliceFloat64 {
	sf64 := SliceFloat64{Elements: []float64{}}
	return sf64
}

func (sf64 *SliceFloat64) AddElement(num float64) {
	sf64.Elements = append(sf64.Elements, num)
}

func (sf64 *SliceFloat64) Count() int {
	return sf64.Len()
}

func (sf64 *SliceFloat64) Len() int {
	return len(sf64.Elements)
}

func (sf64 *SliceFloat64) Min() (float64, error) {
	if len(sf64.Elements) == 0 {
		return 0, errors.New("List is empty")
	}
	if !sort.Float64sAreSorted(sf64.Elements) {
		sort.Float64s(sf64.Elements)
	}
	return sf64.Elements[0], nil
}

func (sf64 *SliceFloat64) Max() (float64, error) {
	if len(sf64.Elements) == 0 {
		return 0, errors.New("List is empty")
	}
	if !sort.Float64sAreSorted(sf64.Elements) {
		sort.Float64s(sf64.Elements)
	}
	return sf64.Elements[len(sf64.Elements)-1], nil
}

func (sf64 *SliceFloat64) Sum() (float64, error) {
	if len(sf64.Elements) == 0 {
		return 0, errors.New("List is empty")
	}
	sum := float64(0)
	for _, num := range sf64.Elements {
		sum += num
	}
	return sum, nil
}

func (sf64 *SliceFloat64) Average() (float64, error) {
	return sf64.Mean()
}

func (sf64 *SliceFloat64) Mean() (float64, error) {
	if len(sf64.Elements) == 0 {
		return 0, errors.New("List is empty")
	}
	sum, err := sf64.Sum()
	if err != nil {
		return 0, err
	}
	return sum / float64(len(sf64.Elements)), nil
}

func (sf64 *SliceFloat64) Median() (float64, error) {
	if len(sf64.Elements) == 0 {
		return 0, errors.New("List is empty")
	}
	if !sort.Float64sAreSorted(sf64.Elements) {
		sort.Float64s(sf64.Elements)
	}
	mid := int64(float64(len(sf64.Elements)) / 2)
	return sf64.Elements[mid], nil
}

func (sf64 *SliceFloat64) Stats() (SliceFloat64Stats, error) {
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
	return stats, nil
}

type SliceFloat64Stats struct {
	Len    int
	Max    float64
	Mean   float64
	Median float64
	Min    float64
	Sum    float64
}

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
