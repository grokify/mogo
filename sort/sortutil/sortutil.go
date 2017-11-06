package sortutil

import (
	"errors"
	"sort"
	"strings"
)

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

func Int64s(a []int64) { sort.Sort(Int64Slice(a)) }

// Int64Slice attaches the methods of Interface to []iint64, sorting in increasing order.
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

func Uint16s(a []uint16) { sort.Sort(Uint16Slice(a)) }
