package stringsutil

import (
	"testing"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var equalFoldFullTests = []struct {
	s      string
	t      string
	sCased string
	tCased string
	caser  cases.Caser
	want   bool
}{
	// grüßen example from: https://stackoverflow.com/q/43059909/1908967
	{"grüßen", "GRÜSSEN", "grüssen", "grüssen", defaultCaser, true},
	{"grüßen", "GRÜSSEN ", "grüssen", "grüssen ", defaultCaser, false},
	{"grüßen", "GRÜßEN ", "grüssen", "grüssen ", defaultCaser, false},
	{"grüßen", "GRÜSSEN", "Grüßen", "GRÜSSEN", cases.Title(language.German, cases.NoLower), false},
	{"grüßen", "Grüßen", "Grüßen", "Grüßen", cases.Title(language.German, cases.NoLower), true},
	{"GRÜSSEN", "GRÜSSEN", "GRÜSSEN", "GRÜSSEN", cases.Title(language.German, cases.NoLower), true},
}

func TestEqualFoldFull(t *testing.T) {
	for i, tt := range equalFoldFullTests {
		gotSCased := tt.caser.String(tt.s)
		if tt.sCased != gotSCased {
			t.Errorf("caser.String(\"%s\") S Test [%d] Error: want [%v], got [%v]",
				tt.s, i, tt.sCased, gotSCased)
		}
		gotTCased := tt.caser.String(tt.t)
		if tt.tCased != gotTCased {
			t.Errorf("caser.String(\"%s\") T Test [%d] Error: want [%v], got [%v]",
				tt.t, i, tt.tCased, gotTCased)
		}
		gotEqualFold := EqualFoldFull(tt.s, tt.t, &tt.caser)
		if tt.want != gotEqualFold {
			t.Errorf("stringsutil.EqualFoldFull(\"%s\", \"%s\") Error: want [%v], got [%v]",
				tt.s, tt.t, tt.want, gotEqualFold)
		}
	}
}
