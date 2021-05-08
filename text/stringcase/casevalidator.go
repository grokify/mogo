package stringcase

import (
	"fmt"
	"regexp"
	"strings"
)

/*

https://stackoverflow.com/questions/1128305/regex-for-pascalcased-words-aka-camelcased-with-leading-uppercase-letter
https://gist.github.com/manjeettahkur/ff114ef92d8ffee1b797091ff77ea89f
https://google.github.io/styleguide/javaguide.html#s5.3-camel-case

*/

const (
	CamelCase  = "camelCase"
	KebabCase  = "kebab-case"
	PascalCase = "PascalCase"
	SnakeCase  = "snake_case"
)

var mapCaseConst = map[string]string{
	"camel":       CamelCase,
	"camelcase":   CamelCase,
	"camel-case":  CamelCase,
	"camel_case":  CamelCase,
	"kebab":       KebabCase,
	"kebabcase":   KebabCase,
	"kebab-case":  KebabCase,
	"kebab_case":  KebabCase,
	"pascal":      PascalCase,
	"pascalcase":  PascalCase,
	"pascal-case": PascalCase,
	"pascal_case": PascalCase,
	"snake":       SnakeCase,
	"snakecase":   SnakeCase,
	"snake-case":  SnakeCase,
	"snake_case":  SnakeCase,
}

func Parse(s string) (string, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if caseConst, ok := mapCaseConst[s]; ok {
		return caseConst, nil
	}
	return "", fmt.Errorf("case [%s] not parsed", s)
}

func IsCase(caseType, s string) (bool, error) {
	caseTypeCanonical, err := Parse(caseType)
	if err != nil {
		return false, err
	}
	switch caseTypeCanonical {
	case CamelCase:
		return IsCamelCase(s), nil
	case KebabCase:
		return IsKebabCase(s), nil
	case PascalCase:
		return IsPascalCase(s), nil
	case SnakeCase:
		return IsSnakeCase(s), nil
	}
	return false, fmt.Errorf("unknown string case type [%s]", caseType)
}

var (
	rxCamelCase         = regexp.MustCompile(`^[a-z][0-9A-Za-z]*$`)
	rxKebabCase         = regexp.MustCompile(`^[a-z][0-9a-z-]*$`)
	rxPascalCase        = regexp.MustCompile(`^[A-Z][0-9A-Za-z]*$`)
	rxSnakeCase         = regexp.MustCompile(`^[a-z][0-9a-z_]*$`)
	rxCamelCaseIdSuffix = regexp.MustCompile(`[0-9a-z](I[dD])$`)
)

// IsCamelCase returns if a string is camelCase or not.
func IsCamelCase(input string) bool {
	if !rxCamelCase.MatchString(input) {
		return false
	}
	m := rxCamelCaseIdSuffix.FindStringSubmatch(input)
	if len(m) == 2 {
		if m[1] != "Id" {
			return false
		}
	}
	return true
}

func IsKebabCase(input string) bool {
	return rxKebabCase.MatchString(input)
}

// IsPascalCase returns if a string is PascalCase or not.
func IsPascalCase(input string) bool {
	if !rxPascalCase.MatchString(input) {
		return false
	}
	m := rxCamelCaseIdSuffix.FindStringSubmatch(input)
	if len(m) == 2 {
		if m[1] != "Id" {
			return false
		}
	}
	return true
}

func IsSnakeCase(input string) bool {
	return rxSnakeCase.MatchString(input)
}

var rxFirstAlphaUpper = regexp.MustCompile(`^[A-Z]`)

// IsFirstAlphaUpper returns if the first character is
// a capital [A-Z] character.
func IsFirstAlphaUpper(s string) bool {
	return rxFirstAlphaUpper.MatchString(s)
}
