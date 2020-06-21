package stringsutil

import (
	"strings"
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
	{"HelloWorld", "helloWorld"},
	{"helloWorld", "helloWorld"}}

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
	{"HelloWorld", "HelloWorld"},
	{"helloWorld", "HelloWorld"}}

func TestToUpperFirst(t *testing.T) {
	for _, tt := range toUpperFirstTests {
		got := ToUpperFirst(tt.v, false)

		if got != tt.want {
			t.Errorf("strutil.ToUpperFirst() Error: with [%v], want [%v], got [%v]",
				tt.v, tt.want, got)
		}
	}
}

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

var toOppositeTests = []struct {
	v    string
	want string
}{
	{lowerUpper, upperLower},
	{"hello", "HELLO"},
	{"Hello", "hELLO"},
	{"HELLO", "hello"},
	{"1Hello", "1hELLO"},
}

func TestToOpposite(t *testing.T) {
	for _, tt := range toOppositeTests {
		got := ToOpposite(tt.v)

		if got != tt.want {
			t.Errorf("strutil.ToOpposite(%v) Error: want [%v], got [%v]",
				tt.v, tt.want, got)
		}
	}
}
