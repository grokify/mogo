package stringsutil

import (
	"testing"
)

var toLowerFirstTests = []struct {
	v    string
	want string
}{
	{"HelloWorld", "helloWorld"},
	{"helloWorld", "helloWorld"}}

func TestToLowerFirst(t *testing.T) {
	for _, tt := range toLowerFirstTests {
		try := tt.v
		got := ToLowerFirst(try)
		if got != tt.want {
			t.Errorf("ToLowerFirst failed")
		}
	}
}

var toUpperFirstTests = []struct {
	v    string
	want string
}{
	{"HelloWorld", "HelloWorld"},
	{"helloWorld", "HelloWorld"}}

func TestToUpperFirst(t *testing.T) {
	for _, tt := range toUpperFirstTests {
		try := tt.v
		got := ToUpperFirst(try)
		if got != tt.want {
			t.Errorf("ToUpperFirst failed")
		}
	}
}
