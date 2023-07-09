package strconvutil

import (
	"sort"
	"strconv"

	"golang.org/x/exp/constraints"
)

// SliceAtoi converts a slice of string integers.
func SliceAtoi(strings []string) ([]int, error) {
	ints := []int{}
	for _, s := range strings {
		thisInt, err := strconv.Atoi(s)
		if err != nil {
			return ints, err
		}
		ints = append(ints, thisInt)
	}
	return ints, nil
}

// SliceAtoiSort converts and sorts a slice of string integers.
func SliceAtoiSort(strings []string) ([]int, error) {
	ints, err := SliceAtoi(strings)
	if err != nil {
		return ints, err
	}
	intSlice := sort.IntSlice(ints)
	intSlice.Sort()
	return intSlice, nil
}

// SliceItoa converts a slice of `constraints.Integer` to a slice of `string`.
func SliceItoa[S ~[]E, E constraints.Integer](elems S) []string {
	strs := []string{}
	for _, v := range elems {
		strs = append(strs, strconv.Itoa(int(v)))
	}
	return strs
}
