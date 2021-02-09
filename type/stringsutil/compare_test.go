package stringsutil

import (
	"testing"
)

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
