package maputil

import "strconv"

type MapCompInt[C comparable] map[C]int

func (mci MapCompInt[C]) ReverseCounts() map[int]int {
	out := map[int]int{}
	for _, v := range mci {
		out[v]++
	}
	return out
}

func (mci MapCompInt[C]) ReverseCountsString() map[string]int {
	out := map[string]int{}
	for _, v := range mci {
		out[strconv.Itoa(v)]++
	}
	return out
}
