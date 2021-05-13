package stringcase

import (
	"regexp"
	"strings"

	"github.com/grokify/simplego/type/stringsutil"
)

func CaseKebabToCamel(s string) string {
	return stringsutil.ToLowerFirst(ToPascalCase(s))
}

func CaseKebabToPascal(s string) string {
	return ToPascalCase(s)
}

func CaseKebabToSnake(s string) string {
	return rxHypen.ReplaceAllString(
		strings.ToLower(strings.TrimSpace(s)), "_")
}

func CaseSnakeToCamel(s string) string {
	return stringsutil.ToLowerFirst(ToPascalCase(s))
}

func CaseSnakeToKebab(s string) string {
	return rxUnderscore.ReplaceAllString(
		strings.ToLower(strings.TrimSpace(s)), "-")
}

func CaseSnakeToPascal(s string) string {
	return ToPascalCase(s)
}

var rxSplitCase = regexp.MustCompile(`[\s_\-;:~]`)

// ToCamelCase converts a string to camel case as `camelCase`.
func ToCamelCase(s string) string {
	return stringsutil.ToLowerFirst(ToPascalCase(s))
}

// ToPascalCase converts a string to Pascal case as `PascalCase`.
func ToPascalCase(s string) string {
	parts := stringsutil.SliceCondenseSpace(rxSplitCase.Split(s, -1), false, false)
	for i, part := range parts {
		parts[i] = stringsutil.ToUpperFirst(part, true)
	}
	return strings.Join(parts, "")
}
