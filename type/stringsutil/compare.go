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
