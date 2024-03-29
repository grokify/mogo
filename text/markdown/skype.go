package markdown

import (
	"fmt"
	"regexp"
)

var (
	rxHTTPURLPrefix *regexp.Regexp = regexp.MustCompile(`^(i?)https?://`)
	rxSkypeLink     *regexp.Regexp = regexp.MustCompile(`<([^><\|]*?)\|([^>]*?)>`)
	rxBacktick3     *regexp.Regexp = regexp.MustCompile(`^\s*` + "```.+```" + `\s*$`)
)

// SkypeToMarkdown converts Skype markup to Markdown. This is specifically
// useful for converting Slack messages to Markdown. The `stripURLAutoLink`
// parameter will remove links when they are within 3 backticks and the
// link innerHTML and URL match. Default is `true`.
func SkypeToMarkdown(input string, stripURLAutoLink bool) string {
	output := input
	backtick3 := rxBacktick3.MatchString(output)
	m := rxSkypeLink.FindAllStringSubmatch(input, -1)
	for _, n := range m {
		rxlink := regexp.MustCompile(regexp.QuoteMeta(n[0]))
		if stripURLAutoLink && backtick3 && n[1] == n[2] && rxHTTPURLPrefix.MatchString(n[1]) {
			output = rxlink.ReplaceAllString(output, n[1])
		} else {
			mkdn := fmt.Sprintf("[%s](%s)", n[2], n[1])
			output = rxlink.ReplaceAllString(output, mkdn)
		}
	}
	return output
}
