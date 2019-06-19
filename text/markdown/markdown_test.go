package markdown

import (
	"testing"
)

var skypeToMarkdownTests = []struct {
	v    string
	want string
}{
	{"```<https://github.com/grokify|https://github.com/grokify>```", "```https://github.com/grokify```"},
	{`<https://github.com/grokify|Grokify>`, `[Grokify](https://github.com/grokify)`},
	{`ABC <https://github.com/grokify|Grokify>`, `ABC [Grokify](https://github.com/grokify)`},
}

func TestSkypeToMarkdown(t *testing.T) {
	for _, tt := range skypeToMarkdownTests {
		got := SkypeToMarkdown(tt.v, true)
		if got != tt.want {
			t.Errorf("markdown.SkypeToMarkdown(\"%v\") Mismatch: want [%v] got [%v]", tt.v, tt.want, got)
		}
	}
}
