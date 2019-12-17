package stringsutil

import "strings"

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
