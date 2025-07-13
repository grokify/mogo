package unicodeutil

import (
	"testing"
)

// TestRemoveDiacritics ensures timeutil.DateDMYHM2 is parsed to GMT timezone.
func TestRemoveDiacritics(t *testing.T) {
	var removeDiacriticsTests = []struct {
		v    string
		vMap map[rune][]rune
		want string
	}{
		{"João", nil, "Joao"},
		{"München", nil, "Munchen"},
		{"item™", map[rune][]rune{'™': {'t', 'm'}}, "itemtm"}, // mkdocs example
	}

	for _, tt := range removeDiacriticsTests {
		try := RemoveDiacritics(tt.v, tt.vMap)
		if try != tt.want {
			t.Errorf("unicodeutil.RemoveDiacritics(\"%s\") Mismatch: want [%s], got [%s]", tt.v, tt.want, try)
		}
	}
}

func TestUnescape(t *testing.T) {
	var unescapeTests = []struct {
		v    string
		want string
	}{
		{`M\u00fcnchen`, "München"},
		{`"M\u00fcnchen"`, `"München"`},
	}

	for _, tt := range unescapeTests {
		unescaped, err := Unescape(tt.v)
		if err != nil {
			t.Errorf("unicodeutil.Unescape(\"%s\") error: (%s)", tt.v, err.Error())
		}
		if unescaped != tt.want {
			t.Errorf("unicodeutil.Unescape(\"%s\") mismatch: want (%s), got (%s)", tt.v, tt.want, unescaped)
		}
	}
}
