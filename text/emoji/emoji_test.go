package emoji

import (
	"testing"
)

var emoji2AsciiTests = []struct {
	v           string
	wantAscii   string
	wantUnicode string
}{
	{`:+1:`, `+1`, `+1`},
	{`:sweat_smile:`, `':)`, `😅`},
	{`:confused: :sweat_smile:`, `>:\ ':)`, `😕 😅`},
}

func TestEmojiToAscii(t *testing.T) {
	conv := NewConverter()
	for _, tt := range emoji2AsciiTests {
		gotAscii := conv.ConvertShortcodesString(tt.v, Ascii)
		if gotAscii != tt.wantAscii {
			t.Errorf("converter.ConvertString(\"%v\", Ascii) Mismatch: want [%v] got [%v]", tt.v, tt.wantAscii, gotAscii)
		}
		gotUnicode := conv.ConvertShortcodesString(tt.v, Unicode)
		if gotUnicode != tt.wantUnicode {
			t.Errorf("converter.ConvertString(\"%v\", Unicode) Mismatch: want [%v] got [%v]", tt.v, tt.wantUnicode, gotUnicode)
		}
	}
}
