package stringsutil

import (
	"strings"
	"testing"
)

var containsTests = []struct {
	haystack string
	needles  []string
	all      bool
	lc       bool
	trim     bool
	want     bool
}{
	{"hello world", []string{"world", "  Hello   "}, false, true, true, true},
	{"hello world", []string{"world", "hello"}, true, true, true, true},
	{"hello world", []string{"worl", "foo"}, false, true, true, true},
	{"hello world", []string{"worl", "foo"}, true, true, true, false},
}

func TestContainsMore(t *testing.T) {
	for _, tt := range containsTests {
		try := ContainsMore(tt.haystack, tt.needles, tt.all, tt.lc, tt.trim)
		if tt.want != try {
			t.Errorf("stringsutil.ContainsMore(\"%s\", \"%s\", %v, %v, %v) Error: want [%v], got [%v]",
				tt.haystack,
				strings.Join(tt.needles, ","), tt.all, tt.lc, tt.trim, tt.want, try)
		}
	}
}

var endsWithTests = []struct {
	fullstring string
	substring  string
	result     bool
}{
	{"hello world", "world", true},
	{"hello world", "worl", false},
}

func TestEndsWith(t *testing.T) {
	for _, tt := range endsWithTests {
		result := EndsWith(tt.fullstring, tt.substring)
		if tt.result != result {
			t.Errorf("stringsutil.EndsWith(\"%s\", \"%s\") Error: want [%v], got [%v]",
				tt.fullstring, tt.substring, tt.result, result)
		}
	}
}
