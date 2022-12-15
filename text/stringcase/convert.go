package stringcase

import (
	"errors"
	"regexp"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
	"github.com/iancoleman/strcase"
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
	return strcase.ToKebab(s)
	// return strings.Join(toParts(strings.ToLower(s)), "-")
}

func ToSnakeCase(s string) string {
	return strcase.ToSnake(s)
	// return strings.Join(toParts(strings.ToLower(s)), "_")
}

func NoOp(s string) string {
	return s
}

func toParts(s string) []string {
	return stringsutil.SliceCondenseSpace(rxSplitCase.Split(s, -1), false, false)
}

func FuncToWantCaseMore(c string, overrides map[string]string) (func(string) string, error) {
	tocase, err := FuncToWantCase(c)
	if err != nil {
		return NoOp, err
	}
	if len(overrides) == 0 {
		return tocase, nil
	}
	return func(s string) string {
		if ovr, ok := overrides[s]; ok {
			return ovr
		}
		return tocase(s)
	}, nil
}

func FuncToWantCase(c string) (func(string) string, error) {
	canonical, err := Parse(c)
	if err != nil {
		return strcase.ToSnake, err
	}
	switch canonical {
	case CamelCase:
		return strcase.ToLowerCamel, nil
	case KebabCase:
		return strcase.ToKebab, nil
	case PascalCase:
		return strcase.ToCamel, nil
	case SnakeCase:
		return strcase.ToSnake, nil
	default:
		return strcase.ToSnake, errors.New("case not supported")
	}
}

/*
func FuncToWantCaseOld(c string) (func(string) string, error) {
	parsed, err := Parse(c)
	if err != nil {
		return NoOp, err
	}
	switch parsed {
	case CamelCase:
		return ToCamelCase, nil
	case KebabCase:
		return ToKebabCase, nil
	case PascalCase:
		return ToPascalCase, nil
	case SnakeCase:
		return ToSnakeCase, nil
	}
	return NoOp, ErrUnknownCaseString // should never hit this.
}
*/

func FuncToWantCaseOrNoOp(c string) func(string) string {
	return FuncToWantCaseOrDefault(c, NoOp)
}

// FuncToWantCaseOrDefault returns an ToWantCase function, or default if case is not preseent or
// parseable. For flexibility, if a `nil` func is passed, `nil`, is returned.
func FuncToWantCaseOrDefault(c string, defaultFunc func(string) string) func(string) string {
	wantFunc, err := FuncToWantCase(c)
	if err != nil {
		return defaultFunc
	}
	return wantFunc
}
