package bigint

import (
	"testing"
)

var encodeStringTests = []struct {
	data []byte
	base int
	enc  string
}{
	{[]byte(`{"foo":"bar"}`), 62, "3iDIA3jnw58n8DhGln"},
	{[]byte(`{"foo":"bar"}`), 36, "q9txxhmgdvb1hq1gvon1"},
	{[]byte(`{"foo":"bar"}`), 16, "7b22666f6f223a22626172227d"},
}

func TestEncodeString(t *testing.T) {
	for _, tt := range encodeStringTests {
		enc, err := EncodeToString(tt.base, tt.data)
		if err != nil {
			t.Errorf("bigutil.EncodeToString(%d,%v):err [%v]", tt.base, string(tt.data), err.Error())
		}
		if tt.enc != enc {
			t.Errorf("bigutil.EncodeToString(%d,%v): want [%v], got [%v]", tt.base, string(tt.data), tt.enc, enc)
		}
		enc2 := MustEncodeToString(tt.base, tt.data)
		if tt.enc != enc2 {
			t.Errorf("bigutil.MustEncodeToString(%d,%v): want [%v], got [%v]", tt.base, string(tt.data), tt.enc, enc)
		}
		dec, err := DecodeString(tt.base, enc)
		if err != nil {
			t.Errorf("bigutil.DecodeString(%d,%v):err [%v]", tt.base, enc, err.Error())
		}
		if string(tt.data) != string(dec) {
			t.Errorf("bigutil.DecodeString(%d,%v): want [%v], got [%v]", tt.base, enc, string(tt.data), string(dec))
		}
	}
}
