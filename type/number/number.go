package number

import "sort"

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
