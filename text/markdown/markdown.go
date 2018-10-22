package markdown

import (
	"fmt"
	"regexp"
)

func BoldText(haystack, needle string) string {
	output := haystack
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
