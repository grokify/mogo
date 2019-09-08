package md5

import (
	"testing"

	"github.com/grokify/gotilla/encoding/base36"
)

var encodeBase36HexStringTests = []struct {
	v    string
	want string
}{
	{"5eb63bbbe01eeed093cb22bb8f5acdc3", "5luw5ld8t195dpiliva0krvsz"},
	{"ed076287532e86365e841e92bfc50d8c", "e16cs890ihyk8hvpfezbncfpo"},
}

func TestEncodeBase36HexString(t *testing.T) {
	for _, tt := range encodeBase36HexStringTests {
		enc := base36.Encode36HexString(tt.v)

		if enc != tt.want {
			t.Errorf("base36.Encode36String(%v): want [%v], got [%v]", tt.v, tt.want, enc)
		}
	}
}

var md5Base36Tests = []struct {
	v      string
	want10 string
	want36 string
}{
	{"hello world", "125893641179230474042701625388361764291", "5luw5ld8t195dpiliva0krvsz"},
	{"Hello World!", "315065379476721403163906509030895717772", "e16cs890ihyk8hvpfezbncfpo"}}

func TestMd5Base36(t *testing.T) {
	for _, tt := range md5Base36Tests {
		enc36 := Md5Base36(tt.v)
		enc10 := Md5Base10(tt.v)

		if enc36 != tt.want36 {
			t.Errorf("base36.Md5Base36(%v): want [%v], got [%v]", tt.v, tt.want36, enc36)
		}
		if enc10 != tt.want10 {
			t.Errorf("base36.Md5Base10(%v): want [%v], got [%v]", tt.v, tt.want10, enc10)
		}
	}
}
