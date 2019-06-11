package emoji

import (
	"testing"
)

var emoji2AsciiTests = []struct {
	v    string
	want string
}{
	{`:+1:`, `+1`},
	{`:sweat_smile:`, `':)`},
	{`:confused: :sweat_smile:`, `>:\ ':)`},
}

func TestEmojiToAscii(t *testing.T) {
	conv := NewConverter()
	for _, tt := range emoji2AsciiTests {
		got := conv.EmojiToAscii(tt.v)
		if got != tt.want {
			t.Errorf("converter.EmojiToAscii(\"%v\") Mismatch: want [%v] got [%v]", tt.v, tt.want, got)
		}
	}
}
