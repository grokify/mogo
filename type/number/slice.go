package number

import (
	"golang.org/x/exp/constraints"

	"github.com/grokify/mogo/strconv/strconvutil"
)

type Integers[C constraints.Integer] []C

func (ints Integers[C]) MinMax() (*C, *C) {
	if len(ints) == 0 {
		return nil, nil
	}
	var min C
	var max C
	for i, v := range ints {
		if i == 0 {
			min = v
			max = v
		} else {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
	}
	return &min, &max
}

func (ints Integers[C]) Sum() C {
	var sum C
	for _, v := range ints {
		sum += v
	}
	return sum
}

func (ints Integers[C]) SumString() string {
	return strconvutil.Itoa(ints.Sum())
}
