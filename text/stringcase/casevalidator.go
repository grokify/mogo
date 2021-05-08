package stringcase

import (
	"fmt"
	"regexp"
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

type CaseValidator struct{}

func IsCase(caseType, s string) (bool, error) {
	switch caseType {
	case CamelCase:
		{
			return IsCamelCase(s), nil
		}
	case PascalCase:
		{
			return IsPascalCase(s), nil
		}
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
