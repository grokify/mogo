package strslices

import (
	"strings"
	"testing"
)

var indexMultiTests = []struct {
	s      string
	substr []string
	idx    int
}{
	{"hello", []string{"foo", "bar", "hello"}, 0},
	{"hello", []string{"foo", "bar", "ello"}, 1},
	{"hello", []string{"foo", "bar", "baz"}, -1},
	{"hello", []string{"foo", "bar", "hello"}, 0},
}

func TestIndexMulti(t *testing.T) {
	for _, tt := range indexMultiTests {
		idx := IndexMulti(tt.s, tt.substr...)
		if idx != tt.idx {
			t.Errorf("stringsutil.IndexMore(%s, []string{%s}) Error: want [%d], got [%d]",
				tt.s, strings.Join(tt.substr, ","), tt.idx, idx)
		}
	}
}
