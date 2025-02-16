package stringsutil

import (
	"strings"

	"golang.org/x/text/cases"
)

// SliceIndex returns the index of the first match using `=`. Returns -1 if not found.
// if `equalFold` is selected and `caser` is `nil`, the default caser will be used.
func SliceIndex(haystack []string, needle string, equalFold bool, caser *cases.Caser) int {
	if len(haystack) == 0 {
		return -1
	}
	for i, hay := range haystack {
		if !equalFold && needle == hay {
			return i
		} else if equalFold && EqualFoldFull(needle, hay, caser) {
			return i
		}
	}
	return -1
}

// SliceIndexContains returns the index of the first match
// using `strings.Contains()`. Returns -1 if not found.
func SliceIndexContains(s []string, substr string) int {
	for i, si := range s {
		if strings.Contains(si, substr) {
			return i
		}
	}
	return -1
}

// SliceIndexMore returns the index of an element in a
// string slice. Returns -1 if not found.
func SliceIndexMore(haystack []string, needle string, trimSpace, toLower bool, matchType MatchType) int {
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
		if matchType == MatchStringSuffix {
			if strings.HasSuffix(needle, hay) {
				return idx
			}
		} else if matchType == MatchStringPrefix {
			if strings.HasPrefix(needle, hay) {
				return idx
			}
		} else {
			if needle == hay {
				return idx
			}
		}
	}
	return -1
}

// SliceIndexOrEmpty returns the element at the index provided or an empty string.
func SliceIndexOrEmpty(s []string, index int) string {
	if index < 0 || index >= len(s) {
		return ""
	} else {
		return s[index]
	}
}

func SliceLineHasIndex(haystack []string, needle string, wantIndex int) bool {
	for _, line := range haystack {
		try := strings.Index(line, needle)
		if try == wantIndex {
			return true
		}
	}
	return false
}

func Slice2FilterLinesHaveIndex(groups [][]string, needle string, wantIndex int) [][]string {
	newGroups := [][]string{}
	for _, grp := range groups {
		if SliceLineHasIndex(grp, needle, wantIndex) {
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
