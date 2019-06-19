package markdown

import (
	"fmt"
	"regexp"
)

var (
	rxHttpUrlPrefix *regexp.Regexp = regexp.MustCompile(`^(i?)https?://`)
	rxSkypeLink     *regexp.Regexp = regexp.MustCompile(`<([^><\|]*?)\|([^>]*?)>`)
	rxBacktick3     *regexp.Regexp = regexp.MustCompile(`^\s*` + "```.+```" + `\s*$`)
)

// BoldText bodifies the identified text. It looks for start of words
// using a word boundary and will arbirarily end to match words with
// different suffixes.
func BoldText(haystack, needle string) string {
	output := haystack
	if len(needle) == 0 {
		return "**" + haystack + "**"
	}
	return regexp.MustCompile(`(?i)\b(`+regexp.QuoteMeta(needle)+`)`).ReplaceAllString(output, "**$1**")
}

func UrlToMarkdownLinkHostname(url string) string {
	rx := regexp.MustCompile(`(?i)^https?://([^/]+)(/[^/])`)
	m := rx.FindStringSubmatch(url)
	if len(m) > 1 {
		suffix := ""
		if len(m) > 2 {
			suffix = "..."
		}
		return fmt.Sprintf("[%v%v](%v)", m[1], suffix, url)
	}
	return url
}

// SkypeToMarkdown converts Skype markup to Markdown. This is specifically
// useful for converting Slack messages to Markdown. The `stripUrlAutoLink`
// parameter will remove links when they are within 3 backticks and the
// link innerHTML and URl match.
func SkypeToMarkdown(input string, stripUrlAutoLink bool) string {
	output := input
	backtick3 := rxBacktick3.MatchString(output)
	m := rxSkypeLink.FindAllStringSubmatch(input, -1)
	for _, n := range m {
		rxlink := regexp.MustCompile(regexp.QuoteMeta(n[0]))
		if stripUrlAutoLink && backtick3 && n[1] == n[2] && rxHttpUrlPrefix.MatchString(n[1]) {
			output = rxlink.ReplaceAllString(output, n[1])
		} else {
			mkdn := fmt.Sprintf("[%s](%s)", n[2], n[1])
			output = rxlink.ReplaceAllString(output, mkdn)
		}
	}
	return output
}
