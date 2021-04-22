package urlutil

import (
	"testing"
)

var uriSchemeTests = []struct {
	v    string
	want string
}{
	{"https://example.com", "https"},
	{"rc6://example", "rc6"}}

func TestUriScheme(t *testing.T) {
	for _, tt := range uriSchemeTests {
		try := tt.v
		got := UriScheme(try)
		if got != tt.want {
			t.Errorf("GetScheme(%v) failed want [%v] got [%v]", tt.v, tt.want, got)
		}
	}
}

var isHttpTests = []struct {
	v         string
	inclHttp  bool
	inclHttps bool
	want      bool
}{
	{"http://example.com", false, false, false},
	{"http://example.com", true, false, true},
	{"http://example.com", false, true, false},
	{"http://example.com", true, true, true},
	{"https://example.com", false, false, false},
	{"https://example.com", true, false, false},
	{"https://example.com", false, true, true},
	{"https://example.com", true, true, true},
	{"email://example.com", true, true, false},
	{"tel://example.com", true, true, false},
}

func TestIsHttp(t *testing.T) {
	for _, tt := range isHttpTests {
		try := tt.v
		got := IsHttp(try, tt.inclHttp, tt.inclHttps)
		if got != tt.want {
			t.Errorf("IsHttp(%v) failed want [%v] got [%v]", tt.v, tt.want, got)
		}
	}
}
