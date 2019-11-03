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
