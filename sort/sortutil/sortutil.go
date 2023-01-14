package sortutil

import (
	"errors"
	"sort"
	"strings"

	"golang.org/x/exp/constraints"
)

// Slice sorts a slice of items that comply wth `constraints.Ordered`.
func Slice[E constraints.Ordered](s []E) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
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
	return "", errors.New("string not found")
}

/*
// Int64Slice attaches the methods of Interface to []int64, sorting in increasing order.
type Int32Slice []int32

func (p Int32Slice) Len() int           { return len(p) }
func (p Int32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p Int32Slice) Sort() { sort.Sort(p) }

// Int64Slice attaches the methods of Interface to []int64, sorting in increasing order.
type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p Int64Slice) Sort() { sort.Sort(p) }

// Uint16Slice attaches the methods of Interface to []uint16, sorting in increasing order.
type Uint16Slice []uint16

func (p Uint16Slice) Len() int           { return len(p) }
func (p Uint16Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint16Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p Uint16Slice) Sort() { sort.Sort(p) }

// Convenience wrappers for common cases

// Int64s sorts a slice of int64s in increasing order.
func Int32s(a []int32) { sort.Sort(Int32Slice(a)) }

// Int64s sorts a slice of int64s in increasing order.
func Int64s(a []int64) { sort.Sort(Int64Slice(a)) }

// Uint16s sorts a slice of uint16s in increasing order.
func Uint16s(a []uint16) { sort.Sort(Uint16Slice(a)) }
*/
