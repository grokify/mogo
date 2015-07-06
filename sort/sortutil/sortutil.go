package sortutil

import (
	"sort"
)

func Int64Sorted(int64s []int64) []int64 {
	ints := Int64sToInts(int64s)
	sort.Ints(ints)
	return IntsToInt64s(ints)
}

func Int64sToInts(int64s []int64) []int {
	ints := []int{}
	for _, x := range int64s {
		ints = append(ints, int(x))
	}
	return ints
}

func IntsToInt64s(ints []int) []int64 {
	int64s := []int64{}
	for _, x := range ints {
		int64s = append(int64s, int64(x))
	}
	return int64s
}
