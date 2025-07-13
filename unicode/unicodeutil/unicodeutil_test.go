package unicodeutil

import (
	"testing"
)

var unescapeTests = []struct {
	v    string
	want string
}{
	{`M\u00fcnchen`, "München"},
	{`"M\u00fcnchen"`, `"München"`},
}

func TestUnescape(t *testing.T) {
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
