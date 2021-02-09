package stringsutil

import (
	"testing"
)

var toCaseTests = []struct {
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

func TestToCase(t *testing.T) {
	for _, tt := range toCaseTests {
		tryCamelCase := ToCamelCase(tt.anyString)
		if tryCamelCase != tt.camelCase {
			t.Errorf("stringsutil.ToCamelCase(\"%s\") Error: want [%s], got [%s]",
				tt.anyString, tt.camelCase, tryCamelCase)
		}
		tryPascalCase := ToPascalCase(tt.anyString)
		if tryPascalCase != tt.pascalCase {
			t.Errorf("stringsutil.ToPascalCase(\"%s\") Error: want [%s], got [%s]",
				tt.anyString, tt.pascalCase, tryPascalCase)
		}
	}
}
