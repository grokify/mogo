package stringsutil

import (
	"regexp"
	"strings"
)

var rxSplitCase = regexp.MustCompile(`[\s_\-;:~]`)

// ToCamelCase converts a string to camel case as `camelCase`.
func ToCamelCase(s string) string {
	return ToLowerFirst(ToPascalCase(s))
}

// ToPascalCase converts a string to Pascal case as `PascalCase`.
func ToPascalCase(s string) string {
	parts := SliceCondenseSpace(rxSplitCase.Split(s, -1), false, false)
	for i, part := range parts {
		parts[i] = ToUpperFirst(part, true)
	}
	return strings.Join(parts, "")
}
