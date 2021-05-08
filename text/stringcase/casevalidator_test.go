package stringcase

import (
	"testing"
)

var caseTests = []struct {
	str      string
	isKebab  bool // lower case
	isSnake  bool // lower case
	isCamel  bool // upper case
	isPascal bool // upper case
}{
	{"camelCase", false, false, true, false},
	{"camelCaseId", false, false, true, false},
	{"camelCaseID", false, false, false, false},
	{"PascalCase", false, false, false, true},
	{"PascalCaseId", false, false, false, true},
	{"PascalCaseID", false, false, false, false},
}

func TestCase(t *testing.T) {
	for _, tt := range caseTests {
		isCamel := IsCamelCase(tt.str)
		if tt.isCamel {
			if !isCamel {
				t.Errorf("stringcase.IsCamelCase(\"%s\") Mismatch: want [%v] got [%v]", tt.str, tt.isCamel, isCamel)
			}
		} else {
			if isCamel {
				t.Errorf("stringcase.IsCamelCase(\"%s\") Mismatch: want [%v] got [%v]", tt.str, tt.isCamel, isCamel)
			}
		}

		isKebab := IsKebabCase(tt.str)
		if tt.isKebab {
			if !isKebab {
				t.Errorf("stringcase.IsKebabCase(\"%s\") Mismatch: want [%v] got [%v]", tt.str, tt.isKebab, isKebab)
			}
		} else {
			if isKebab {
				t.Errorf("stringcase.IsKebabCase(\"%s\") Mismatch: want [%v] got [%v]", tt.str, tt.isKebab, isKebab)
			}
		}

		isPascal := IsPascalCase(tt.str)
		if tt.isPascal {
			if !isPascal {
				t.Errorf("stringcase.IsPascalCase(\"%s\") Mismatch: want [%v] got [%v]", tt.str, tt.isPascal, isPascal)
			}
		} else {
			if isPascal {
				t.Errorf("stringcase.IsPascalCase(\"%s\") Mismatch: want [%v] got [%v]", tt.str, tt.isPascal, isPascal)
			}
		}

		isSnake := IsSnakeCase(tt.str)
		if tt.isSnake {
			if !isSnake {
				t.Errorf("stringcase.IsSnakeCase(\"%s\") Mismatch: want [%v] got [%v]", tt.str, tt.isSnake, isSnake)
			}
		} else {
			if isSnake {
				t.Errorf("stringcase.IsSnakeCase(\"%s\") Mismatch: want [%v] got [%v]", tt.str, tt.isSnake, isSnake)
			}
		}

	}
}
