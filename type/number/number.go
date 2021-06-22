package number

import (
	"sort"
	"strconv"
)

func SliceInt64ToFloat64(src []int64) []float64 {
	out := []float64{}
	for _, in := range src {
		out = append(out, float64(in))
	}
	return out
}

func IntSliceDedupe(elems []int, sortElems bool) []int {
	if len(elems) == 0 {
		return elems
	}
	deduped := []int{}
	mapInts := map[int]int{}
	for _, el := range elems {
		if _, ok := mapInts[el]; !ok {
			deduped = append(deduped, el)
			mapInts[el] = 1
		}
	}
	if sortElems {
		sort.Ints(deduped)
	}
	return deduped
}

func IntLength(num int) uint {
	if num < 0 {
		num = num * -1
	}
	return uint(len(strconv.Itoa(num)))
}
