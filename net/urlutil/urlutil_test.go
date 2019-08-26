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
	{" ---   hello World 中文---   ", "hello-world-中文"}}

func TestToSlugLowerString(t *testing.T) {
	for _, tt := range toSlugLowerStringTests {
		try := tt.v
		got := ToSlugLowerString(try)
		if got != tt.want {
			t.Errorf("ToSlugLowerString failed want [%v] got [%v]", tt.want, got)
		}
	}
}

var uriCondenseTests = []struct {
	v    string
	want string
}{
	{"https://abc//def//", "https://abc/def/"},
	{"  https://abc//def//  ", "https://abc/def/"},
	{"https://////abc///def/", "https://abc/def/"}}

func TestUriCondense(t *testing.T) {
	for _, tt := range uriCondenseTests {
		try := tt.v
		got := UriCondense(try)
		if got != tt.want {
			t.Errorf("UriCondense(%v) failed want [%v] got [%v]", tt.v, tt.want, got)
		}
	}
}
