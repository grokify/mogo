package stringcase

import (
	"testing"
)

var convTests = []struct {
	camel  string
	kebab  string
	pascal string
	snake  string
}{
	{"helloWorld", "hello-world", "HelloWorld", "hello_world"},
	{"fooBarBazQux", "foo-bar-baz-qux", "FooBarBazQux", "foo_bar_baz_qux"},
}

func TestConv(t *testing.T) {
	for _, tt := range convTests {

		tryCamel := CaseKebabToCamel(tt.kebab)
		if tryCamel != tt.camel {
			t.Errorf("stringcase.CaseKebabToCamel(\"%s\") Mismatch: want [%v] got [%v]",
				tt.kebab, tt.camel, tryCamel)
		}

		tryPascal := CaseKebabToPascal(tt.kebab)
		if tryPascal != tt.pascal {
			t.Errorf("stringcase.CaseKebabToPascal(\"%s\") Mismatch: want [%v] got [%v]",
				tt.kebab, tt.pascal, tryPascal)
		}

		trySnake := CaseKebabToSnake(tt.kebab)
		if trySnake != tt.snake {
			t.Errorf("stringcase.CaseKebabToSnake(\"%s\") Mismatch: want [%v] got [%v]",
				tt.kebab, tt.snake, trySnake)
		}

		tryCamel = CaseSnakeToCamel(tt.snake)
		if tryCamel != tt.camel {
			t.Errorf("stringcase.CaseSnakeToCamel(\"%s\") Mismatch: want [%v] got [%v]",
				tt.snake, tt.camel, tryCamel)
		}

		tryKebab := CaseSnakeToKebab(tt.snake)
		if tryKebab != tt.kebab {
			t.Errorf("stringcase.CaseSnakeToKebab(\"%s\") Mismatch: want [%v] got [%v]",
				tt.snake, tt.kebab, tryKebab)
		}

		tryPascal = CaseSnakeToPascal(tt.snake)
		if tryPascal != tt.pascal {
			t.Errorf("stringcase.CaseSnakeToPascal(\"%s\") Mismatch: want [%v] got [%v]",
				tt.snake, tt.pascal, tryPascal)
		}

	}
}

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
