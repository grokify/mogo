package htmlutil

import (
	"testing"
)

var htmlToTextTests = []struct {
	v    string
	want string
}{
	{"<p>Foo</p><p>Bar</p>", "Foo\n\nBar"},
}

func TestHTMLToText(t *testing.T) {
	for _, tt := range htmlToTextTests {
		got := HtmlToText(tt.v)

		if got != tt.want {
			t.Errorf("htmlutil.TestHTMLToText(`%s`) Error: want [%v], got [%v]",
				tt.v, tt.want, got)
		}
	}
}
