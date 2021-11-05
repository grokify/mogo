package markdown

import "testing"

var linkifyTests = []struct {
	url      string
	text     string
	markdown string
}{
	{
		"https://example.com", "", "[https://example.com](https://example.com)"},
	{
		"  https://example.com  ", "", "[https://example.com](https://example.com)"},
	{
		"https://example.com", "Example", "[Example](https://example.com)"},
	{
		"  https://example.com", "Example", "[Example](https://example.com)"},
	{
		"", "Example", "Example"},
	{
		"   ", "Example", "Example"},
}

func TestLinkify(t *testing.T) {
	for _, tt := range linkifyTests {
		got := Linkify(tt.url, tt.text)
		if got != tt.markdown {
			t.Errorf("markdown.Linkify(\"%s\", \"%s\") Mismatch: want [%s] got [%s]",
				tt.url, tt.text, tt.markdown, got)
		}
	}
}
