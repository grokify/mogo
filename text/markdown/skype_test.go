package markdown

import (
	"testing"
)

var skypeToMarkdownTests = []struct {
	skype    string
	markdown string
}{
	{
		"```<https://github.com/grokify|https://github.com/grokify>```",
		"```https://github.com/grokify```"},
	{
		`<https://github.com/grokify|Grokify>`,
		`[Grokify](https://github.com/grokify)`},
	{
		`ABC <https://github.com/grokify|Grokify>`,
		`ABC [Grokify](https://github.com/grokify)`},
	{
		`<https://gitlab.com/example/example-project/-/commit/8b2c5d00884b08d12a9d8f97e3ce81de30d78b8b|8b2c5d00>: Add LICENSE - Grokify`,
		`[8b2c5d00](https://gitlab.com/example/example-project/-/commit/8b2c5d00884b08d12a9d8f97e3ce81de30d78b8b): Add LICENSE - Grokify`},
	{
		`John Wang pushed to branch <https://gitlab.com/grokify/example-project/commits/main|main> of <https://gitlab.com/grokify/example-project|John Wang / Example Project> (<https://gitlab.com/grokify/example-project/compare/90c523f219e3d41b098c74678ecc50ac581428b6...8b2c5d00884b08d12a9d8f97e3ce81de30d78b8b|Compare changes>)\n[8b2c5d00](https://gitlab.com/grokify/example-project/-/commit/8b2c5d00884b08d12a9d8f97e3ce81de30d78b8b): Add LICENSE - John Wang`,
		`John Wang pushed to branch [main](https://gitlab.com/grokify/example-project/commits/main) of [John Wang / Example Project](https://gitlab.com/grokify/example-project) ([Compare changes](https://gitlab.com/grokify/example-project/compare/90c523f219e3d41b098c74678ecc50ac581428b6...8b2c5d00884b08d12a9d8f97e3ce81de30d78b8b))\n[8b2c5d00](https://gitlab.com/grokify/example-project/-/commit/8b2c5d00884b08d12a9d8f97e3ce81de30d78b8b): Add LICENSE - John Wang`},
}

func TestSkypeToMarkdown(t *testing.T) {
	for _, tt := range skypeToMarkdownTests {
		got := SkypeToMarkdown(tt.skype, true)
		if got != tt.markdown {
			t.Errorf("markdown.SkypeToMarkdown(\"%v\") Mismatch: want [%v] got [%v]",
				tt.skype, tt.markdown, got)
		}
	}
}
