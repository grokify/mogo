package stringsutil

import "strings"

// SliceIndex returns the index of the first match
// using `=`. Returns -1 if not found.
func SliceIndex(haystack []string, needle string) int {
	if len(haystack) == 0 {
		return -1
	}
	for i, el := range haystack {
		if needle == el {
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
		if matchType == MatchHasSuffix {
			if strings.HasSuffix(needle, hay) {
				return idx
			}
		} else if matchType == MatchHasPrefix {
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

// SliceIndexOrEmpty returns the element at the index
// provided or an empty string.
func SliceIndexOrEmpty(s []string, index uint64) string {
	if int(index) >= len(s) {
		return ""
	}
	return s[index]
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
