package stringsutil

import (
	"regexp"
	"strings"
)

var rxSplitCase = regexp.MustCompile(`[\s_\-;:~]`)

// ToCamelCase converts a string to camel case as `camelCase`.
func ToCamelCase(s string) string {
	parts := SliceCondenseSpace(rxSplitCase.Split(s, -1), false, false)
	for i, part := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(part)
		} else {
			parts[i] = ToUpperFirst(part, true)
		}
	}
	return strings.Join(parts, "")
}
