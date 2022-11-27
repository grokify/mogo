package stringsutil

import (
	"testing"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var equalFoldFullTests = []struct {
	s     string
	t     string
	caser cases.Caser
	want  bool
}{
	{"grüßen", "GRÜSSEN", defaultCaser, true},
	{"grüßen", "GRÜSSEN ", defaultCaser, false},
	{"grüßen", "GRÜßEN ", defaultCaser, false},
	{"grüßen", "GRÜSSEN", cases.Title(language.German, cases.NoLower), false},
}

func TestEqualFoldFull(t *testing.T) {
	for _, tt := range equalFoldFullTests {
		got := EqualFoldFull(tt.s, tt.t, &tt.caser)
		if tt.want != got {
			t.Errorf("stringsutil.EqualFoldFull(\"%s\", \"%s\") Error: want [%v], got [%v]",
				tt.s, tt.t, tt.want, got)
		}
	}
}
