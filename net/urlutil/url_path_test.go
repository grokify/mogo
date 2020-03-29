package urlutil

import (
	"testing"
)

var pathLeafTests = []struct {
	v    string
	want string
}{
	{"https://example.com/v5", "v5"},
	{"https://customer.example.com:0/{v5}", "{v5}"},
	{"https://foobar.com/v8.0", "v8.0"},
	{"https://foobar/v1/v2/v8.0/", "v8.0"},
}

func TestGetPathLeaf(t *testing.T) {
	for _, tt := range pathLeafTests {
		u1, err := GetPathLeaf(tt.v)
		if err != nil {
			t.Errorf("urlutil.GetPathLeaf() Error: input [%v], want [%v], got error [%v]",
				tt.v, tt.want, err.Error())
		}
		if u1 != tt.want {
			t.Errorf("urlutil.GetPathLeaf() Failure: input [%v], want [%v], got [%v]",
				tt.v, tt.want, u1)
		}
	}
}

var modifyPathTests = []struct {
	rawurl  string
	newpath string
	want    string
}{
	{"https://example.com/foo/", "/", "https://example.com"},
	{"https://example.com/foo", "bar", "https://example.com/bar"},
	{"https://example.com/foo/bar/baz", "bar", "https://example.com/bar"},
	{"https://example.com/foo/bar/baz?foo1=bar1#abc", "bar", "https://example.com/bar?foo1=bar1#abc"},
}

func TestModifyPath(t *testing.T) {
	for _, tt := range modifyPathTests {
		u1, err := ModifyPath(tt.rawurl, tt.newpath)
		if err != nil {
			t.Errorf("urlutil.ModifyPath() Error: input [%v, %v], want [%v], got error [%v]",
				tt.rawurl, tt.newpath, tt.want, err.Error())
		}
		if u1 != tt.want {
			t.Errorf("urlutil.ModifyPath() Failure: input [%v, %v], want [%v], got [%v]",
				tt.rawurl, tt.newpath, tt.want, u1)
		}
	}
}
