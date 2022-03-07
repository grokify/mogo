package emoji

import (
	"testing"
)

var emoji2ASCIITests = []struct {
	v           string
	wantASCII   string
	wantUnicode string
}{
	{`:-1:`, `-1`, `👎`},
	{`:+1:`, `+1`, `👍`},
	{`:sweat_smile:`, `':)`, `😅`},
	{`:confused: :sweat_smile:`, `>:\ ':)`, `😕 😅`},
}

func TestEmojiToASCII(t *testing.T) {
	conv := NewConverter()
	for _, tt := range emoji2ASCIITests {
		gotASCII := conv.ConvertShortcodesString(tt.v, ASCII)
		if gotASCII != tt.wantASCII {
			t.Errorf("converter.ConvertString(\"%v\", ASCII) Mismatch: want [%v] got [%v]", tt.v, tt.wantASCII, gotASCII)
		}
		gotUnicode := conv.ConvertShortcodesString(tt.v, Unicode)
		if gotUnicode != tt.wantUnicode {
			t.Errorf("converter.ConvertString(\"%v\", Unicode) Mismatch: want [%v] got [%v]", tt.v, tt.wantUnicode, gotUnicode)
		}
	}
}
