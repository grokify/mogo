package textutil

import (
	"testing"
)

var removeDiacriticsTests = []struct {
	v    string
	want string
}{
	{"å", "a"},
	{"ß", "ss"},
	{"Jesús", "Jesus"},
	{"žůžo", "zuzo"},
}

func TestRemoveDiacritics(t *testing.T) {
	for _, tt := range removeDiacriticsTests {
		try, err := RemoveDiacritics(tt.v)
		if err != nil {
			t.Errorf("strconvutil.RemoveDiacritics(\"%s\") Error: [%s]",
				tt.v, err.Error())
		}
		if err == nil && try != tt.want {
			t.Errorf("strconvutil.RemoveDiacritics(\"%s\" Error: want [%s], got [%s]",
				tt.v, tt.want, try)
		}
	}
}
