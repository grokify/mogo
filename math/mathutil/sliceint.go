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

func (sint *SliceInt) AddInt(num int) {
	sint.Elements = append(sint.Elements, num)
}

func (sint *SliceInt) AddInt32(num int32) {
	sint.Elements = append(sint.Elements, int(num))
}

func (sint *SliceInt) AddInt64(num int64) {
	sint.Elements = append(sint.Elements, int(num))
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
