package maputil

import "strconv"

type MapComparableInt[C comparable] map[C]int

func (mci MapComparableInt[C]) ReverseCounts() map[int]int {
	out := map[int]int{}
	for _, v := range mci {
		out[v]++
	}
	return out
}

func (mci MapComparableInt[C]) ReverseCountsString() map[string]int {
	out := map[string]int{}
	for _, v := range mci {
		out[strconv.Itoa(v)]++
	}
	return out
}
