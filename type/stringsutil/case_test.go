package stringsutil

import (
	"testing"
)

var equalFoldFullTests = []struct {
	s    string
	t    string
	want bool
}{
	{"grüßen", "GRÜSSEN", true},
	{"grüßen", "GRÜSSEN ", false},
	{"grüßen", "GRÜßEN ", false},
}

func TestEqualFoldFull(t *testing.T) {
	for _, tt := range equalFoldFullTests {
		got := EqualFoldFull(tt.s, tt.t)
		if tt.want != got {
			t.Errorf("stringsutil.EqualFoldFull(\"%s\", \"%s\") Error: want [%v], got [%v]",
				tt.s, tt.t, tt.want, got)
		}
	}
}
