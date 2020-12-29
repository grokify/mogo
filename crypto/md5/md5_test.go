package md5

import (
	"testing"

	"github.com/grokify/simplego/encoding/base36"
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
	v       string
	want10  string
	want36  string
	want62  string
	want62u string
}{
	{"hello world", "125893641179230474042701625388361764291",
		"5luw5ld8t195dpiliva0krvsz", "2SIyH7gjExZ74B2pirixcT", "2siYh7GJeXz74b2PIRIXCt"},
	{"Hello World!", "315065379476721403163906509030895717772",
		"e16cs890ihyk8hvpfezbncfpo", "7dgyuMkhmWALzZmAxQB3Y0", "7DGYUmKHMwalZzMaXqb3y0"}}

func TestMd5Base36(t *testing.T) {
	for _, tt := range md5Base36Tests {
		enc62 := Md5Base62(tt.v)
		enc62u := Md5Base62UpperFirst(tt.v)
		enc36 := Md5Base36(tt.v)
		enc10 := Md5Base10(tt.v)

		if enc62 != tt.want62 {
			t.Errorf("md5.Md5Base62(%v): want [%v], got [%v]", tt.v, tt.want62, enc62)
		}

		if enc62u != tt.want62u {
			t.Errorf("md5.Md5Base62UpperFirst(%v): want [%v], got [%v]", tt.v, tt.want62u, enc62u)
		}
		if enc36 != tt.want36 {
			t.Errorf("md5.Md5Base36(%v): want [%v], got [%v]", tt.v, tt.want36, enc36)
		}
		if enc10 != tt.want10 {
			t.Errorf("md5.Md5Base10(%v): want [%v], got [%v]", tt.v, tt.want10, enc10)
		}
	}
}
