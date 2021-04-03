package jsonutil

import (
	"testing"
)

var jsonschemaEscapeTests = []struct {
	unescaped string
	escaped   string
}{
	{"a~b", "a~0b"},
	{"a/b", "a~1b"},
	{"a~/~/b", "a~0~1~0~1b"},
	{"a~~~///b", "a~0~0~0~1~1~1b"},
}

func TestJSONSchemaEscape(t *testing.T) {
	for _, tt := range jsonschemaEscapeTests {
		tryEscape := PropertyNameEscape(tt.unescaped)
		if tryEscape != tt.escaped {
			t.Errorf("jsonutil.PropertyNameEscape: want [%v] got [%v]", tt.escaped, tryEscape)
		}
		tryUnescape := PropertyNameUnescape(tt.escaped)
		if tryUnescape != tt.unescaped {
			t.Errorf("jsonutil.PropertyNameUnescape: want [%v] got [%v]", tt.unescaped, tryUnescape)
		}
	}
}
