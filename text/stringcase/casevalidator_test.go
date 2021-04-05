package stringcase

import (
	"testing"
)

var caseTests = []struct {
	str      string
	isCamel  bool
	isPascal bool
}{
	{"camelCase", true, false},
	{"camelCaseId", true, false},
	{"camelCaseID", false, false},
	{"PascalCase", false, true},
	{"PamelCaseId", false, true},
	{"PamelCaseID", false, false},
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
	}
}
