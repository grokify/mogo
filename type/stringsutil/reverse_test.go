package stringsutil

import (
	"testing"
)

var reverseTests = []struct {
	Str          string
	Reverse      string
	Substr       string
	ReverseIndex int
}{
	{"Hello World", "dlroW olleH", " World", 0},
	{" Hello World ", " dlroW olleH ", " World ", 0},
	{" Hello World2 ", " 2dlroW olleH ", " World2", 1},
	{" Hello World3 ", " 3dlroW olleH ", "World4 ", -1},
}

func TestReverse(t *testing.T) {
	for _, tt := range reverseTests {
		rev := Reverse(tt.Str)
		if rev != tt.Reverse {
			t.Errorf("stringsutil.Reverse(\"%s\") want [%s] got [%s]",
				tt.Str, tt.Reverse, rev)
		}
		rIdx := ReverseIndex(tt.Str, tt.Substr)
		if rIdx != tt.ReverseIndex {
			t.Errorf("stringsutil.ReverseIndex(\"%s\", \"%s\") want [%d] got [%d]",
				tt.Str, tt.Substr, tt.ReverseIndex, rIdx)
		}
	}
}
