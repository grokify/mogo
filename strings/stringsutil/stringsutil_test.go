package stringsutil

import (
	"strings"
	"testing"
)

var joinLiteraryTests = []struct {
	v        []string
	sep      string
	joinWord string
	want     string
}{
	{[]string{}, ",", "and", ""},
	{[]string{"Foo"}, ",", "and", "Foo"},
	{[]string{"Foo", "Bar"}, ",", "and", "Foo and Bar"},
	{[]string{"Foo", "Bar", "Baz"}, ",", "and", "Foo, Bar, and Baz"},
	{[]string{"Foo", "Bar", "Bax", "Qux"}, ",", "and", "Foo, Bar, Bax, and Qux"}}

func TestJoinLiterary(t *testing.T) {
	for _, tt := range joinLiteraryTests {
		try := tt.v
		got := JoinLiterary(try, tt.sep, tt.joinWord)
		if got != tt.want {
			t.Errorf("TestJoinLanguage failed: Have [%v] Got [%v] Want [%v]",
				strings.Join(tt.v, ", "),
				got,
				tt.want)
		}
	}
}

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
