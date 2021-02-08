package stringsutil

import (
	"testing"
)

var toCamelCaseTests = []struct {
	anyString  string
	camelCase  string
	pascalCase string
}{
	{"Lorem ipsum dolor sit amet", "loremIpsumDolorSitAmet", "LoremIpsumDolorSitAmet"},
	{"Lorem IPSUM", "loremIpsum", "LoremIpsum"},
	{"Lorem Ipsum", "loremIpsum", "LoremIpsum"},
	{"snake_case", "snakeCase", "SnakeCase"},
	{"kebab-case", "kebabCase", "KebabCase"},
}

func TestToCamelCase(t *testing.T) {
	for _, tt := range toCamelCaseTests {
		tryCamelCase := ToCamelCase(tt.anyString)
		if tryCamelCase != tt.camelCase {
			t.Errorf("stringsutil.ToCamelCase(\"%s\") Error: want [%s], got [%s]",
				tt.anyString, tt.camelCase, tryCamelCase)
		}
	}
}
