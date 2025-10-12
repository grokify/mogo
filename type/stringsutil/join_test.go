package stringsutil

import (
	"strings"
	"testing"
)

func TestJoinLiterary(t *testing.T) {
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
