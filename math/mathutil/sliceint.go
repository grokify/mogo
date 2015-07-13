package mathutil

import (
	"errors"
	"sort"
)

type SliceInt struct {
	Elements []int
}

func NewSliceInt() SliceInt {
	sint := SliceInt{Elements: []int{}}
	return sint
}

func (sint *SliceInt) AddElement(num int) {
	sint.Elements = append(sint.Elements, num)
}

func (sint *SliceInt) Count() int {
	return sint.Len()
}

func (sint *SliceInt) Len() int {
	return len(sint.Elements)
}

func (sint *SliceInt) Min() (int, error) {
	if len(sint.Elements) == 0 {
		return 0, errors.New("List is empty")
	}
	if !sort.IntsAreSorted(sint.Elements) {
		sort.Ints(sint.Elements)
	}
	return sint.Elements[0], nil
}

func (sint *SliceInt) Max() (int, error) {
	if len(sint.Elements) == 0 {
		return 0, errors.New("List is empty")
	}
	if !sort.IntsAreSorted(sint.Elements) {
		sort.Ints(sint.Elements)
	}
	return sint.Elements[len(sint.Elements)-1], nil
}

func (sint *SliceInt) Sum() (int, error) {
	if len(sint.Elements) == 0 {
		return 0, errors.New("List is empty")
	}
	sum := int(0)
	for _, num := range sint.Elements {
		sum += num
	}
	return sum, nil
}

func (sint *SliceInt) Average() (float64, error) {
	return sint.Mean()
}

func (sint *SliceInt) Mean() (float64, error) {
	if len(sint.Elements) == 0 {
		return 0, errors.New("List is empty")
	}
	sum, err := sint.Sum()
	if err != nil {
		return 0, err
	}
	return float64(sum) / float64(len(sint.Elements)), nil
}

func (sint *SliceInt) Median() (int, error) {
	if len(sint.Elements) == 0 {
		return 0, errors.New("List is empty")
	}
	if !sort.IntsAreSorted(sint.Elements) {
		sort.Ints(sint.Elements)
	}
	mid := int64(float64(len(sint.Elements)) / 2)
	return sint.Elements[mid], nil
}

func (sint *SliceInt) Stats() SliceIntStats {
	stats := NewSliceIntStats()
	stats.Len = sint.Len()
	max, err := sint.Max()
	if err == nil {
		stats.Max = max
	}
	min, err := sint.Min()
	if err == nil {
		stats.Min = min
	}
	mean, err := sint.Mean()
	if err == nil {
		stats.Mean = mean
	}
	median, err := sint.Median()
	if err == nil {
		stats.Median = median
	}
	sum, err := sint.Sum()
	if err == nil {
		stats.Sum = sum
	}
	return stats
}

type SliceIntStats struct {
	Len    int
	Max    int
	Mean   float64
	Median int
	Min    int
	Sum    int
}

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
