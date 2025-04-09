package strslices

import "strings"

// Contains checks a slice of string to validate if individual
// string elements contain the supplied match string. A `1` is
// returned if all slice strings match, a `-1` is returned if no
// slices strings match, and a `0` is returned if more than one,
// but less than all slice strings match. An empty slice returns
// a `-1` as there are no members to match.
func Contains(s []string, v string) int { // all, some, none
	if len(s) == 0 {
		return -1
	}
	haveContains := false
	haveNotContains := false
	for _, si := range s {
		if strings.Contains(si, v) {
			haveContains = true
		} else {
			haveNotContains = true
		}
	}
	if haveContains && !haveNotContains {
		return 1 // all have
	} else if haveNotContains && !haveContains {
		return -1 // none have
	} else {
		return 0 // some have
	}
}
