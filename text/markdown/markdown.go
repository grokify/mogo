package markdown

import (
	"fmt"
	"regexp"
	"strings"
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

func URLToMarkdownLinkHostname(url string) string {
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

// Linkify constructs a link from url and text inputs.
// This function does not handle escaping so it should
// be done before hand, e.g. using `\[` instead of `[`
// by itself in the `text`.
func Linkify(url, text string) string {
	url = strings.TrimSpace(url)
	if len(url) == 0 {
		return text
	}
	if len(text) == 0 {
		return "[" + url + "](" + url + ")"
	}
	return "[" + text + "](" + url + ")"
}

var rxParseLink = regexp.MustCompile(`\[(.*)\]\((.*)\)`)

// ParseLink returns text and url.
func ParseLink(s string) (string, string) {
	m := rxParseLink.FindStringSubmatch(s)
	if len(m) > 0 {
		return m[1], m[2]
	}
	return s, ""
}
