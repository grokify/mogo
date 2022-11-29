package jsonpointer

import (
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

var jsonPointerTests = []struct {
	jsonPointer string
	document    string
	path        []string
}{
	{"/components/schemas/FooBar", "", []string{"components", "schemas", "FooBar"}},
	{"#/components/schemas/FooBar", "", []string{"components", "schemas", "FooBar"}},
	{"mydoc.yaml#/components/schemas/FooBar", "mydoc.yaml", []string{"components", "schemas", "FooBar"}},
}

// TestParseJSONPointer ensures the `ParseJSONPointer` is working properly.
func TestParseJSONPointer(t *testing.T) {
	for _, tt := range jsonPointerTests {
		ptr, err := ParseJSONPointer(tt.jsonPointer)
		if err != nil {
			t.Errorf("openapi3.ParseJSONPointer(\"%s\") Error [%s]",
				tt.jsonPointer, err.Error())
		}
		if ptr.Document != tt.document {
			t.Errorf("JSONPointer.Document Mismatch: want [%v], got [%v]",
				tt.document, ptr.Document)
		}
		if !slices.Equal(ptr.Path, tt.path) {
			t.Errorf("JSONPointer.Path Mismatch: want [%v], got [%v]",
				strings.Join(tt.path, ", "), strings.Join(ptr.Path, ", "))
		}
	}
}
