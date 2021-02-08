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

func EndsWith(s, substr string) bool {
	idx := strings.Index(s, substr)
	if len(s) >= len(substr) &&
		idx > -1 &&
		idx == len(s)-len(substr) {
		return true
	}
	return false
}
