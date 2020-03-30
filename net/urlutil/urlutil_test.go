package urlutil

import (
	"net/url"
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

var condenseUriTests = []struct {
	v    string
	want string
}{
	{"https://abc//def//", "https://abc/def/"},
	{"https:/abc//def//", "https://abc/def/"},
	{"  https://abc//def//  ", "https://abc/def/"},
	{"https://////abc///def/", "https://abc/def/"}}

func TestCondenseUri(t *testing.T) {
	for _, tt := range condenseUriTests {
		try := tt.v
		got := CondenseUri(try)
		if got != tt.want {
			t.Errorf("UriCondense(%v) failed want [%v] got [%v]", tt.v, tt.want, got)
		}
	}
}

var urlValuesEncodeSortedTests = []struct {
	raw      string
	priority []string
	want     string
}{
	{"foo=baz&foo=bar&foo=quuz", []string{}, "foo=bar&foo=baz&foo=quuz"},
	{"foo=bar&baz=qux&quux=quuz", []string{}, "baz=qux&foo=bar&quux=quuz"},
	{"foo=bar&baz=qux&quux=quuz", []string{"quux"}, "quux=quuz&baz=qux&foo=bar"},
	{"foo=bar&baz=qux&quux=quuz", []string{"quux", "foo", "quux"}, "quux=quuz&foo=bar&baz=qux"},
	{"foo=bar&baz=qux&quux=quuz", []string{"quux", "foo"}, "quux=quuz&foo=bar&baz=qux"},
}

func TestUrlValuesEncodeSorted(t *testing.T) {
	for _, tt := range urlValuesEncodeSortedTests {
		try, err := url.ParseQuery(tt.raw)
		if err != nil {
			t.Errorf("UrlValuesEncodeSorted(...) failed ParseQuery(%v)", tt.raw)
		}
		got := UrlValuesEncodeSorted(try, tt.priority)
		if got != tt.want {
			t.Errorf("UrlValuesEncodeSorted [%v][%v] failed want [%v] got [%v]", tt.raw, tt.priority, tt.want, got)
		}
	}
}

var joinTests = []struct {
	v    []string
	want string
}{
	{[]string{"foo", "bar"}, "foo/bar"},
	{[]string{"foo/", "/bar"}, "foo/bar"},
	{[]string{"foo/bar", "baz"}, "foo/bar/baz"},
	{[]string{"foo", "", "bar"}, "foo/bar"}}

func TestJoin(t *testing.T) {
	for _, tt := range joinTests {
		try := tt.v
		got := Join(try...)
		if got != tt.want {
			t.Errorf("Join failed want [%v] got [%v]", tt.want, got)
		}
	}
}
