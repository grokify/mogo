package strconvutil

import (
	"sort"
	"strconv"
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

func SliceItoa(ints []int) []string {
	strs := []string{}
	for _, intVal := range ints {
		strs = append(strs, strconv.Itoa(intVal))
	}
	return strs
}
