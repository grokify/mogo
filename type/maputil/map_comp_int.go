package maputil

import "strconv"

type MapCompInt[C comparable] map[C]int

func (mci MapCompInt[C]) Percentages() map[C]float32 {
	out := map[C]float32{}
	sum := 0
	for _, v := range mci {
		sum += v
	}
	for k, v := range mci {
		out[k] = float32(v) / float32(sum)
	}
	return out
}

func (mci MapCompInt[C]) PercentagesInt() map[C]int {
	out := map[C]int{}
	sum := 0
	for _, v := range mci {
		sum += v
	}
	for k, v := range mci {
		out[k] = int(100 * (float32(v) / float32(sum)))
	}
	return out
}

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
