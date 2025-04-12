package stringsutil

import (
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

var splitTrimSpaceTests = []struct {
	v         string
	sep       string
	exclEmpty bool
	want      []string
}{
	{" foo , bar,   , baz , qux ", ",", true, []string{"foo", "bar", "baz", "qux"}},
	{" foo , bar,   , baz , qux ", ",", false, []string{"foo", "bar", "", "baz", "qux"}},
}

func TestSplitTrimSpace(t *testing.T) {
	for _, tt := range splitTrimSpaceTests {
		got := SplitTrimSpace(tt.v, tt.sep, tt.exclEmpty)
		if !slices.Equal(tt.want, got) {
			t.Errorf("stringsutil.TestSplitTrimSpace(\"%s\", \"%s\", %v) want (%s) got (%s)",
				tt.v, tt.sep, tt.exclEmpty, strings.Join(tt.want, ","), strings.Join(got, ","))
		}
	}
}
