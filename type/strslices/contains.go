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

func MatchAny(s1, s2 []string, caseInsensive, trimSpace bool) bool {
	for _, s1x := range s1 {
		if trimSpace {
			s1x = strings.TrimSpace(s1x)
		}
		if caseInsensive {
			s1x = strings.ToLower(s1x)
		}
		for _, s2x := range s2 {
			if trimSpace {
				s2x = strings.TrimSpace(s2x)
			}
			if caseInsensive {
				s2x = strings.ToLower(s2x)
			}
			if s1x == s2x {
				return true
			}
		}
	}
	return false
}
