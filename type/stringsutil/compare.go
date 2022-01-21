package stringsutil

import (
	"strings"
)

func Equal(str1, str2 string, trim, lower bool) bool {
	if trim {
		str1 = strings.TrimSpace(str1)
		str2 = strings.TrimSpace(str2)
	}
	if lower {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}
	if str1 == str2 {
		return true
	}
	return false
}

func EndsWith(s string, substrs ...string) bool {
	for _, substr := range substrs {
		idx := strings.Index(s, substr)
		if len(s) >= len(substr) &&
			idx > -1 &&
			idx == len(s)-len(substr) {
			return true
		}
	}
	return false
}

func ContainsMore(s string, substrs []string, all, lc, trimSpaceSubstr bool) bool {
	if lc {
		s = strings.ToLower(s)
	}
	matchAll := true
	for _, sub := range substrs {
		if lc {
			sub = strings.ToLower(sub)
		}
		if trimSpaceSubstr {
			sub = strings.TrimSpace(sub)
		}
		ok := strings.Contains(s, sub)
		if ok && !all {
			return true
		} else if !ok {
			if all {
				return false
			}
			matchAll = false
		}
	}
	return matchAll
}
