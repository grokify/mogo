package slicesutil

import (
	"regexp"
	"sort"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// AppendBulk appends multiple slice to a slice, and skips 0 length slices.
func AppendBulk[S ~[]E, E any](s S, ss []S) S {
	out := slices.Clone(s)
	for _, si := range ss {
		if len(si) > 0 {
			out = append(out, si...)
		}
	}
	return out
}

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

// LengthCounts returns a `map[uint]uint` where the keys are element lengths and the values
// are counts of slices with those lengths.
func LengthCounts[E any](s [][]E) map[uint]uint {
	stats := map[uint]uint{}
	for _, si := range s {
		stats[uint(len(si))]++
	}
	return stats
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

func Prepend[S ~[]E, E any](s []E, e E) []E {
	return append([]E{e}, s...)
}

// Reverse reverses the order of a slice.
func Reverse[E comparable](s []E) {
	// sourced from Stack Overflow under MIT license: https://stackoverflow.com/a/71904070/1908967
	sort.SliceStable(s, func(i, j int) bool {
		return i > j
	})
}

func Shift[S ~[]E, E any](s S) (E, S) {
	if len(s) == 0 {
		return *new(E), []E{}
	}
	return s[0], s[1:]
}

func Sort[E constraints.Ordered](s []E) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

func SortSliceOfSlice[S ~[][]E, E constraints.Ordered | string](s S, indexes ...uint) {
	for _, idx := range indexes {
		sort.Slice(s, func(i, j int) bool {
			return s[i][idx] < s[j][idx]
		})
	}
}

// SplitMaxLength returns a slice of slices where each sub-slice has the max length supplied.
// A supplied `maxLength` of `0` indicates no max length.
func SplitMaxLength[S ~[]E, E any](s S, maxLen int) []S {
	if maxLen <= 0 || len(s) <= maxLen {
		return []S{append(S{}, s...)}
	}
	var split []S
	new := S{}
	for _, e := range s {
		new = append(new, e)
		if len(new) >= maxLen {
			split = append(split, new)
			new = S{}
		}
	}
	if len(new) > 0 {
		split = append(split, new)
	}
	return split
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

func Unique[S ~[]E, E comparable](s S) bool {
	m := map[E]int{}
	for _, e := range s {
		m[e]++
	}
	return len(s) == len(m)
}

// NewWithDefault creates a slice of length `size` which values populated by default value `d`.
func NewWithDefault[E any](size int, d E) []E {
	var s []E
	if size <= 0 {
		return s
	}
	for i := 0; i < size; i++ {
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

func MakeRepeatingElement[V any](l int, v V) []V {
	var out []V
	for i := 0; i < l; i++ {
		out = append(out, v)
	}
	return out
}

func MatrixGetOneOrDefault[C comparable](m [][]C, keyIdx int, keyValue C, wantIdx int, defaultValue C) C {
	if keyIdx < 0 || wantIdx < 0 {
		return defaultValue
	}
	for _, row := range m {
		if keyIdx >= len(row) {
			continue
		}
		if keyValue == row[keyIdx] {
			if wantIdx > len(row) {
				continue
			}
			return row[wantIdx]
		}
	}
	return defaultValue
}

// Split will split a slice into a slice of slices where each slice has a max size `n`.
func Split[S ~[]E, E comparable](s S, n int) []S {
	max := n
	if max <= 0 || max >= len(s) {
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
