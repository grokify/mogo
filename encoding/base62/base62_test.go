package base62

import (
	"testing"
)

var encodeBase36StringTests = []struct {
	plaintext string
	encoded62 string
}{
	{`{"foo":"bar","baz":1}`, "eyJmb17iOiJiYXIiLCJiYXoiOjF8"},
	{"Hello", "SGVsbG7+"},
	{"Hello World", "SGVsbG7gV18ybGQ+"}}

func TestEncodeBase62String(t *testing.T) {
	levels := []int{9}
	for _, tt := range encodeBase36StringTests {
		for _, level := range levels {
			enc := EncodeGzip([]byte(tt.plaintext), level)

			if level == 0 && enc != tt.encoded62 {
				t.Errorf("base62.EncodeGzip(%v): want [%v], got [%v]", tt.plaintext, tt.encoded62, enc)
			}

			enc = StripPadding(enc)

			if !ValidBase62(enc) {
				t.Errorf("base62.EncodeGzip(%v): got [%v], err [%v]", tt.plaintext, enc, "E_NOT_BASE62")
			}

			dec, err := DecodeGunzip(enc)
			if err != nil {
				t.Errorf("base62.DecodeGuzip(%v): want [%v], err [%v]", enc, tt.plaintext, err.Error())
			}

			if string(dec) != tt.plaintext {
				t.Errorf("base62.DecodeGuzip(%v): want [%v], err [%v]", enc, tt.plaintext, string(dec))
			}
		}
	}
}
