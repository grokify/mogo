package stringsutil

import (
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

var sliceTests = []struct {
	v        []string
	want     []string
	condense bool
}{
	{[]string{"  foo  ", " ", "bar"}, []string{"foo", "bar"}, true},
	{[]string{"  foo  ", " ", "bar"}, []string{"foo", "", "bar"}, false},
}

func TestSlices(t *testing.T) {
	for _, tt := range sliceTests {
		got := SliceTrim(tt.v, " ", tt.condense)
		if !slices.Equal(tt.want, got) {
			t.Errorf("stringsutil.SliceTrim(\"%s\") want (%s) got (%s)",
				strings.Join(tt.v, ","), strings.Join(tt.want, ","), strings.Join(got, ","))
		}
	}
}
