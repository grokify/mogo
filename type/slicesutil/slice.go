package slicesutil

import (
	"regexp"
	"sort"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Dedupe returns a string slice with duplicate values removed. First observance is kept.
func Dedupe[S ~[]E, E comparable](s S) S {
	deduped := []E{}
	seen := map[E]int{}
	for _, val := range s {
		if _, ok := seen[val]; ok {
			continue
		}
		seen[val] = 1
		deduped = append(deduped, val)
	}
	return deduped
}

func ElementCounts[E comparable](s []E) map[E]int {
	m := map[E]int{}
	for _, si := range s {
		m[si]++
	}
	return m
}

func MatchFilters[E comparable](s, inclFilters, exclFilters []E, inclAll bool) bool {
	if len(inclFilters) == 0 && len(exclFilters) == 0 {
		return true
	}
	if len(inclFilters) > 0 {
		matches := 0
		for _, nf := range inclFilters {
			idx := slices.Index(s, nf)
			if idx >= 0 {
				matches++
			} else if idx < 0 && inclAll {
				return false
			}
		}
		if matches == 0 {
			return false
		}
	}
	if len(exclFilters) > 0 {
		for _, xf := range exclFilters {
			idx := slices.Index(s, xf)
			if idx >= 0 {
				return false
			}
		}
	}
	return true
}

// Reverse reverses the order of a slice.
func Reverse[E comparable](s []E) {
	// sourced from Stack Overflow under MIT license: https://stackoverflow.com/a/71904070/1908967
	sort.SliceStable(s, func(i, j int) bool {
		return i > j
	})
}

func SortSliceOfSlice[S ~[][]E, E constraints.Ordered | string](s S, indexes ...uint) {
	for _, idx := range indexes {
		sort.Slice(s, func(i, j int) bool {
			return s[i][idx] < s[j][idx]
		})
	}
}

// Sub returns a string slice with duplicate values removed. First observance is kept.
func Sub[S ~[]E, E comparable](s, t S) S {
	filtered := S{}
	toRemove := map[E]int{}
	for _, e := range t {
		toRemove[e]++
	}
	for _, val := range s {
		if _, ok := toRemove[val]; ok {
			continue
		}
		filtered = append(filtered, val)
	}
	return filtered
}

// Sub returns a string slice with duplicate values removed. First observance is kept.
func SubRegexpString(s []string, r *regexp.Regexp) []string {
	filtered := []string{}
	for _, val := range s {
		if r.MatchString(val) {
			continue
		}
		filtered = append(filtered, val)
	}
	return filtered
}

func UniqueValues[S ~[]E, E comparable](s S) bool {
	m := map[E]int{}
	for _, e := range s {
		m[e]++
	}
	return len(s) == len(m)
}

// NewWithDefault creates a slice of length `size` which values populated by default value `d`.
func NewWithDefault[E any](size uint, d E) []E {
	var s []E
	sz := int(size)
	for i := 0; i < sz; i++ {
		s = append(s, d)
	}
	return s
}

// MakeMatrix2D returns a 2-dimensional matrix. Usage as follows
// `a := Make2D[uint8](dy, dx)`.
// Sourced from: https://stackoverflow.com/a/71781206/1908967
func MakeMatrix2D[E any](n, m int) [][]E {
	matrix := make([][]E, n)
	rows := make([]E, n*m)
	for i, startRow := 0, 0; i < n; i, startRow = i+1, startRow+m {
		endRow := startRow + m
		matrix[i] = rows[startRow:endRow:endRow]
	}
	return matrix
}

func MatrixGetOneOrDefault[C comparable](m [][]C, keyIdx uint, keyValue C, wantIdx uint, defaultValue C) C {
	for _, row := range m {
		if int(keyIdx) >= len(row) {
			continue
		}
		if keyValue == row[keyIdx] {
			if int(wantIdx) > len(row) {
				continue
			}
			return row[wantIdx]
		}
	}
	return defaultValue
}

// Split will split a slice into a slice of slices where each slice ha a max size `n`.
func Split[S ~[]E, E comparable](s S, n uint) []S {
	max := int(n)
	if max == 0 || max >= len(s) {
		return []S{s}
	}
	var sos []S
	for _, e := range s {
		if len(sos) == 0 || len(sos[len(sos)-1]) == max {
			sos = append(sos, S{})
		}
		sos[len(sos)-1] = append(sos[len(sos)-1], e)
	}
	return sos
}

func ToMatrix[S ~[]E, E comparable](s S) []S {
	var m []S
	for _, e := range s {
		m = append(m, []E{e})
	}
	return m
}
