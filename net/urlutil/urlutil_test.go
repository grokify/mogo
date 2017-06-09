package urlutil

import (
	"testing"
)

var toSlugLowerStringTests = []struct {
	v    string
	want string
}{
	{"HelloWorld", "helloworld"},
	{"  hello World  ", "hello-world"},
	{"---   hello World 中文---   ", "hello-world-中文"}}

func TestToSlugLowerString(t *testing.T) {
	for _, tt := range toSlugLowerStringTests {
		try := tt.v
		got := ToSlugLowerString(try)
		if got != tt.want {
			t.Errorf("ToSlugLowerString failed want [%v] got [%v]", tt.want, got)
		}
	}
}
