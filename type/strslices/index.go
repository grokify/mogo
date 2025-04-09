package strslices

import (
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
	"golang.org/x/text/cases"
)

// Index returns the index of the first match using `=`. Returns -1 if not found.
// if `equalFold` is selected and `caser` is `nil`, the default caser will be used.
func Index(haystack []string, needle string, equalFold bool, caser *cases.Caser) int {
	if len(haystack) == 0 {
		return -1
	}
	for i, hay := range haystack {
		if !equalFold && needle == hay {
			return i
		} else if equalFold && stringsutil.EqualFoldFull(needle, hay, caser) {
			return i
		}
	}
	return -1
}

// ContainsIndex returns the index of the first match
// using `strings.Contains()`. Returns -1 if not found.
func ContainsIndex(s []string, substr string) int {
	for i, si := range s {
		if strings.Contains(si, substr) {
			return i
		}
	}
	return -1
}

// IndexMore returns the index of an element in a string slice. Returns -1 if not found.
func IndexMore(haystack []string, needle string, trimSpace, toLower bool, matchType stringsutil.MatchType) int {
	if trimSpace {
		needle = strings.TrimSpace(needle)
	}
	if toLower {
		needle = strings.ToLower(needle)
	}
	for idx, hay := range haystack {
		if trimSpace {
			hay = strings.TrimSpace(hay)
		}
		if toLower {
			hay = strings.ToLower(hay)
		}
		switch matchType {
		case stringsutil.MatchStringSuffix:
			if strings.HasSuffix(needle, hay) {
				return idx
			}
		case stringsutil.MatchStringPrefix:
			if strings.HasPrefix(needle, hay) {
				return idx
			}
		default:
			if needle == hay {
				return idx
			}
		}
	}
	return -1
}

// IndexValueOrDefault returns the element at the index provided or the default string.
func IndexValueOrDefault(s []string, index int, def string) string {
	if index < 0 || index >= len(s) {
		return def
	} else {
		return s[index]
	}
}

func ElementHasIndex(haystack []string, needle string, wantIndex int) bool {
	for _, line := range haystack {
		try := strings.Index(line, needle)
		if try == wantIndex {
			return true
		}
	}
	return false
}

func SoSFilterLinesHaveIndex(groups [][]string, needle string, wantIndex int) [][]string {
	newGroups := [][]string{}
	for _, grp := range groups {
		if ElementHasIndex(grp, needle, wantIndex) {
			newGroups = append(newGroups, grp)
		}
	}
	return newGroups
}

// IndexMulti returns the earliest match.
func IndexMulti(s string, substr ...string) int {
	idxm := -1
	for _, sub := range substr {
		idx := strings.Index(s, sub)
		if idx > -1 && (idxm < 0 || idx < idxm) {
			idxm = idx
		}
	}
	return idxm
}
