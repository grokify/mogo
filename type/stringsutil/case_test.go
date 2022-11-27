package stringsutil

import (
	"testing"

	"golang.org/x/text/cases"
)

var equalFoldFullTests = []struct {
	s    string
	t    string
	opts []cases.Option
	want bool
}{
	{"grüßen", "GRÜSSEN", []cases.Option{}, true},
	{"grüßen", "GRÜSSEN ", []cases.Option{}, false},
	{"grüßen", "GRÜßEN ", []cases.Option{}, false},
}

func TestEqualFoldFull(t *testing.T) {
	for _, tt := range equalFoldFullTests {
		got := EqualFoldFull(tt.s, tt.t, tt.opts...)
		if tt.want != got {
			t.Errorf("stringsutil.EqualFoldFull(\"%s\", \"%s\") Error: want [%v], got [%v]",
				tt.s, tt.t, tt.want, got)
		}
		if len(tt.opts) == 0 {
			got := EqualFoldFull(tt.s, tt.t)
			if tt.want != got {
				t.Errorf("stringsutil.EqualFoldFull(\"%s\", \"%s\") Error: want [%v], got [%v]",
					tt.s, tt.t, tt.want, got)
			}
		}
	}
}
