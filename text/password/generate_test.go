package password

import (
	"testing"

	"github.com/grokify/mogo/fmt/fmtutil"
)

var buildAlphabetTests = []struct {
	v    GenerateOpts
	want string
}{
	{GenerateOpts{
		InclLower:     true,
		InclUpper:     true,
		InclNumbers:   true,
		InclSymbols:   true,
		InclAmbiguous: true,
		ExclSimilar:   false,
	}, "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%&*?{}[]()/\\'\"`,;:.<>"},
	{GenerateOpts{
		InclLower:     true,
		InclUpper:     true,
		InclNumbers:   true,
		InclSymbols:   true,
		InclAmbiguous: true,
		ExclSimilar:   true,
	}, "abcdefghjklmnpqrstuvwxyz23456789ABCDEFGHJKMNPQRSTUVWXYZ!@#$%&*?{}[]()/\\'\"`,;:.<>"},
}

func TestBuildAlphabet(t *testing.T) {
	for _, tt := range buildAlphabetTests {
		got := BuildAlphabet(tt.v)
		if got != tt.want {
			fmtutil.PrintJSON(tt.v)
			t.Errorf("password.BuildAlphabet() Mismatch: want (%s) got (%s)",
				tt.want, got)
		}
	}
}
