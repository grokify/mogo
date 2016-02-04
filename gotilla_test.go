package gotilla

import (
	"testing"
)

var bootstrapTests = []struct {
	v    string
	want string
}{
	{"hello", "hello"},
	{"world", "world"}}

func TestBootstrap(t *testing.T) {
	for _, tt := range bootstrapTests {
		got := tt.v
		if got != tt.want {
			t.Errorf("Bootstrap failed")
		}
	}
}
