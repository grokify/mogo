package jsonutil

import (
	"testing"
)

var jsonHashTests = []struct {
	v       map[string]any
	padding rune
	want    string
}{
	{map[string]any{"foo": "bar"}, rune(0), "4EXQNYHIMEW4LRMDGXXZYWQGHXIQF7ILMNWX5M2AOJTHT26OHYIQ"},
	{map[string]any{"foo": 123}, '=', "EF7MB6CKRBML5YFXLMDZQ3HNEWAX6COQNL43CVHRSPRYQD4KOWKQ===="},
}

func TestJSONHashes(t *testing.T) {
	for _, tt := range jsonHashTests {
		try, err := SHA512d256Base32(tt.v, tt.padding)
		if err != nil {
			t.Errorf("jsonutil.SHA512d256Base32: err (%s)", err.Error())
		} else if try != tt.want {
			t.Errorf("jsonutil.SHA512d256Base32(\"%s\", '%v'): want (%s) got (%s)", tt.v, tt.padding, tt.want, try)
		}
	}
}
