package stringcase

import (
	"testing"
)

var convTests = []struct {
	v      string
	camel  string
	kebab  string
	pascal string
	snake  string
}{
	{"hello: World", "helloWorld", "hello-world", "HelloWorld", "hello_world"},
	{"hello: World Rúnar", "helloWorldRunar", "hello-world-runar", "HelloWorldRunar", "hello_world_runar"},
	{"helloWorld", "helloWorld", "hello-world", "HelloWorld", "hello_world"},
	{"fooBarBazQux", "fooBarBazQux", "foo-bar-baz-qux", "FooBarBazQux", "foo_bar_baz_qux"},
}

func TestConv(t *testing.T) {
	for _, tt := range convTests {
		tryKebabRaw := ToKebabCase(tt.v)
		if tryKebabRaw != tt.kebab {
			t.Errorf("stringcase.CaseSnakeToKebab(\"%s\") Mismatch: want [%v] got [%v]",
				tt.snake, tt.kebab, tryKebabRaw)
		}

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
	kebabCase  string
	pascalCase string
	snakeCase  string
}{
	{"Lorem ipsum dolor sit amet", "loremIpsumDolorSitAmet", "lorem-ipsum-dolor-sit-amet", "LoremIpsumDolorSitAmet", "lorem_ipsum_dolor_sit_amet"},
	{"Lorem IPSUM", "loremIpsum", "lorem-ipsum", "LoremIpsum", "lorem_ipsum"},
	{"Lorem Ipsum", "loremIpsum", "lorem-ipsum", "LoremIpsum", "lorem_ipsum"},
	{"snake_case", "snakeCase", "snake-case", "SnakeCase", "snake_case"},
	{"kebab-case", "kebabCase", "kebab-case", "KebabCase", "kebab_case"},
}

func TestToCase(t *testing.T) {
	for _, tt := range toCaseTests {
		tryCamelCase := ToCamelCase(tt.anyString)
		if tryCamelCase != tt.camelCase {
			t.Errorf("stringsutil.ToCamelCase(\"%s\") Error: want [%s], got [%s]",
				tt.anyString, tt.camelCase, tryCamelCase)
		}
		tryKebabCase := ToKebabCase(tt.anyString)
		if tryKebabCase != tt.kebabCase {
			t.Errorf("stringsutil.ToKebabCase(\"%s\") Error: want [%s], got [%s]",
				tt.anyString, tt.kebabCase, tryKebabCase)
		}
		tryPascalCase := ToPascalCase(tt.anyString)
		if tryPascalCase != tt.pascalCase {
			t.Errorf("stringsutil.ToPascalCase(\"%s\") Error: want [%s], got [%s]",
				tt.anyString, tt.pascalCase, tryPascalCase)
		}
		trySnakeCase := ToSnakeCase(tt.anyString)
		if trySnakeCase != tt.snakeCase {
			t.Errorf("stringsutil.ToSnakeCase(\"%s\") Error: want [%s], got [%s]",
				tt.anyString, tt.snakeCase, trySnakeCase)
		}
	}
}
