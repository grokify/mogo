package languageutil

import (
	"strings"
	"testing"

	"golang.org/x/text/language"
)

var joinLanguageTests = []struct {
	v        []string
	sep      string
	joinLang language.Tag
	want     string
}{
	{[]string{}, ",", language.English, ""},
	{[]string{"Foo"}, ",", language.English, "Foo"},
	{[]string{"Foo", "Bar"}, ",", language.English, "Foo and Bar"},
	{[]string{"Foo", "Bar", "Baz"}, ",", language.English, "Foo, Bar, and Baz"},
	{[]string{"Foo", "Bar", "Bax", "Qux"}, ",", language.English, "Foo, Bar, Bax, and Qux"}}

func TestJoinLanguage(t *testing.T) {
	for _, tt := range joinLanguageTests {
		try := tt.v
		got, err := JoinLanguage(try, tt.sep, tt.joinLang)
		if err != nil {
			t.Errorf("languageutil.TestJoinLanguage failed: Have [%v] Got [%v] Want [%v] Err[%v]",
				strings.Join(tt.v, ", "),
				got,
				tt.want,
				err.Error())
		} else if got != tt.want {
			t.Errorf("languageutil.TestJoinLanguage failed: Have [%v] Got [%v] Want [%v]",
				strings.Join(tt.v, ", "),
				got,
				tt.want)
		}
	}
}
