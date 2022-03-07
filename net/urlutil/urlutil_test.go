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
			t.Errorf("func ToSlugLowerString failed want [%v] got [%v]", tt.want, got)
		}
	}
}

var condenseURITests = []struct {
	v    string
	want string
}{
	{"https://abc//def//", "https://abc/def/"},
	{"https:/abc//def//", "https://abc/def/"},
	{"  https://abc//def//  ", "https://abc/def/"},
	{"https://////abc///def/", "https://abc/def/"}}

func TestCondenseURI(t *testing.T) {
	for _, tt := range condenseURITests {
		try := tt.v
		got := CondenseURI(try)
		if got != tt.want {
			t.Errorf("func CondenseURI(%v) failed want [%v] got [%v]", tt.v, tt.want, got)
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
			t.Errorf("test UrlValuesEncodeSorted(...) failed ParseQuery(%v)", tt.raw)
		}
		got := URLValuesEncodeSorted(try, tt.priority)
		if got != tt.want {
			t.Errorf("test UrlValuesEncodeSorted [%v][%v] failed want [%v] got [%v]", tt.raw, tt.priority, tt.want, got)
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

var qsAddTests = []struct {
	baseURL string
	qryKey  string
	qryVal  string
	wantURL string
}{
	{"http://example.com", "foo", "bar", "http://example.com?foo=bar"}}

func TestURLAddQueryString(t *testing.T) {
	for _, tt := range qsAddTests {
		qsMap := map[string][]string{
			tt.qryKey: {tt.qryVal}}
		qsVal := url.Values{}
		qsVal.Set(tt.qryKey, tt.qryVal)

		goURL1, err := URLAddQueryValuesString(tt.baseURL, qsVal)
		if err != nil {
			t.Errorf("Got error [%s]", err.Error())
		}
		if goURL1.String() != tt.wantURL {
			t.Errorf("func URLAddQueryValuesString failed want [%v] got [%v]",
				tt.wantURL, goURL1.String())
		}

		goURL2, err := URLAddQueryString(tt.baseURL, qsMap)
		if err != nil {
			t.Errorf("Got error [%s]", err.Error())
		}
		if goURL2.String() != tt.wantURL {
			t.Errorf("func URLAddQueryString failed want [%v] got [%v]",
				tt.wantURL, goURL2.String())
		}

		goURLInput, err := url.Parse(tt.baseURL)
		if err != nil {
			t.Errorf("Got error url.Parse error [%s]", err.Error())
		}
		goURL3 := URLAddQuery(goURLInput, qsMap)
		if goURL3.String() != tt.wantURL {
			t.Errorf("func URLAddQuery failed want [%v] got [%v]",
				tt.wantURL, goURL3.String())
		}
		goURL4 := URLAddQueryValues(goURLInput, qsVal)
		if goURL4.String() != tt.wantURL {
			t.Errorf("func URLAddQueryValues failed want [%v] got [%v]",
				tt.wantURL, goURL4.String())
		}
	}
}
