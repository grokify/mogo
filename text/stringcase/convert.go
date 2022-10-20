package stringcase

import (
	"regexp"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
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
	parts := toParts(s)
	for i, part := range parts {
		parts[i] = stringsutil.ToUpperFirst(part, true)
	}
	return strings.Join(parts, "")
}

func ToKebabCase(s string) string {
	return strings.Join(toParts(strings.ToLower(s)), "-")
}

func ToSnakeCase(s string) string {
	return strings.Join(toParts(strings.ToLower(s)), "_")
}

func toParts(s string) []string {
	return stringsutil.SliceCondenseSpace(rxSplitCase.Split(s, -1), false, false)
}
