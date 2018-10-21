package strutil

import (
	"testing"
)

var capitalizeTests = []struct {
	v    string
	want string
}{
	{"hello", "Hello"},
	{"Hello", "Hello"},
	{"HELLO", "Hello"},
	{"1Hello", "1hello"},
}

func TestCapitalize(t *testing.T) {
	for _, tt := range capitalizeTests {
		got := Capitalize(tt.v)

		if got != tt.want {
			t.Errorf("strutil.Capitalize() Error: with [%v], want [%v], got [%v]",
				tt.v, tt.want, got)
		}
	}
}

var toLowerFirstTests = []struct {
	v    string
	want string
}{
	{"hello", "hello"},
	{"Hello", "hello"},
	{"1Hello", "1Hello"},
}

func TestToLowerFirst(t *testing.T) {
	for _, tt := range toLowerFirstTests {
		got := ToLowerFirst(tt.v)

		if got != tt.want {
			t.Errorf("strutil.ToLowerFirst() Error: with [%v], want [%v], got [%v]",
				tt.v, tt.want, got)
		}
	}
}

var toUpperFirstTests = []struct {
	v    string
	want string
}{
	{"hello", "Hello"},
	{"Hello", "Hello"},
	{"1Hello", "1Hello"},
}

func TestToUpperFirst(t *testing.T) {
	for _, tt := range toUpperFirstTests {
		got := ToUpperFirst(tt.v)

		if got != tt.want {
			t.Errorf("strutil.ToUpperFirst() Error: with [%v], want [%v], got [%v]",
				tt.v, tt.want, got)
		}
	}
}
