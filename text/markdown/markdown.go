package markdown

import (
	"fmt"
	"regexp"
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
// useful for converting Slack messages to Markdown.
func SkypeToMarkdown(input string) string {
	output := input
	rx := regexp.MustCompile(`<([^><\|]*?)\|([^>]*?)>`)
	m := rx.FindAllStringSubmatch(input, -1)
	for _, n := range m {
		mkdn := fmt.Sprintf("[%s](%s)", n[2], n[1])
		rxlink := regexp.MustCompile(regexp.QuoteMeta(n[0]))
		output = rxlink.ReplaceAllString(output, mkdn)
	}
	return output
}
