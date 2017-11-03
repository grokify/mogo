package sortutil

import (
	"errors"
	"sort"
	"strings"
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

// For now, use only for slices < 100 in length for performance.
// To do: more scalable implementation that uses sorting/searching.
func InArrayStringCaseInsensitive(haystack []string, needle string) (string, error) {
	needleLower := strings.ToLower(strings.TrimSpace(needle))
	for _, canonical := range haystack {
		canonicalLower := strings.ToLower(canonical)
		if canonicalLower == needleLower {
			return canonical, nil
		}
	}
	return "", errors.New("String not found")
}

// Uint16Slice attaches the methods of Interface to []uint16, sorting in increasing order.
type Uint16Slice []uint16

func (p Uint16Slice) Len() int           { return len(p) }
func (p Uint16Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint16Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p Uint16Slice) Sort() { sort.Sort(p) }

func Uint16s(a []uint16) { sort.Sort(Uint16Slice(a)) }
