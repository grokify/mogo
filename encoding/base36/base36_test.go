package base36

import (
	"testing"
)

var encodeBase36StringTests = []struct {
	v    string
	want string
}{
	{"Hello", "3yud78mn"}}

func TestEncodeBase36String(t *testing.T) {
	for _, tt := range encodeBase36StringTests {
		enc := Encode36String(tt.v)

		if enc != tt.want {
			t.Errorf("base36.Encode36String(%v): want [%v], got [%v]", tt.v, tt.want, enc)
		}
	}
}

func TestDecodeBase36String(t *testing.T) {
	for _, tt := range encodeBase36StringTests {
		dec, err := Decode36String(tt.want)
		if err != nil {
			t.Errorf("base36.Decode36String(%v) error [%v]", tt.want, err)
		} else if string(dec) != tt.v {
			t.Errorf("base36.Decode36String(%v): want [%v], got [%v]", tt.v, tt.want, dec)
		}
	}
}

var encodeBase36HexStringTests = []struct {
	v    string
	want string
}{
	{"5eb63bbbe01eeed093cb22bb8f5acdc3", "5luw5ld8t195dpiliva0krvsz"},
	{"ed076287532e86365e841e92bfc50d8c", "e16cs890ihyk8hvpfezbncfpo"},
}

func TestEncodeBase36HexString(t *testing.T) {
	for _, tt := range encodeBase36HexStringTests {
		enc := Encode36HexString(tt.v)

		if enc != tt.want {
			t.Errorf("base36.Encode36String(%v): want [%v], got [%v]", tt.v, tt.want, enc)
		}
	}
}

var md5Base36Tests = []struct {
	v    string
	want string
}{
	{"hello world", "5luw5ld8t195dpiliva0krvsz"},
	{"Hello World!", "e16cs890ihyk8hvpfezbncfpo"}}

func TestMd5Base36(t *testing.T) {
	for _, tt := range md5Base36Tests {
		enc := Md5Base36(tt.v)

		if enc != tt.want {
			t.Errorf("base36.Md5Base36(%v): want [%v], got [%v]", tt.v, tt.want, enc)
		}
	}
}
