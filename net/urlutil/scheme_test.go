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

func TestURIScheme(t *testing.T) {
	for _, tt := range uriSchemeTests {
		try := tt.v
		got := URIScheme(try)
		if got != tt.want {
			t.Errorf("GetScheme(%v) failed want [%v] got [%v]", tt.v, tt.want, got)
		}
	}
}

var isHTTPTests = []struct {
	v         string
	inclHTTP  bool
	inclHTTPS bool
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

func TestIsHTTP(t *testing.T) {
	for _, tt := range isHTTPTests {
		try := tt.v
		got := IsHTTP(try, tt.inclHTTP, tt.inclHTTPS)
		if got != tt.want {
			t.Errorf("func IsHttp(%v) failed want [%v] got [%v]", tt.v, tt.want, got)
		}
	}
}
