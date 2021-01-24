package stringsutil

import (
	"testing"
)

var toCamelCaseTests = []struct {
	anyString string
	camelCase string
}{
	{"Lorem ipsum dolor sit amet", "loremIpsumDolorSitAmet"},
	{"Lorem IPSUM", "loremIpsum"},
	{"Lorem Ipsum", "loremIpsum"},
	{"snake_case", "snakeCase"},
	{"snake-case", "snakeCase"},
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
