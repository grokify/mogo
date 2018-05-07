package stringsutil

import (
	"fmt"
	"regexp"
)

func UrlToMarkdownLinkHostname(url string) string {
	rx := regexp.MustCompile(`(?i)^https?://([^/]+)(/[^/])`)
	m := rx.FindStringSubmatch(url)
	if len(m) > 1 {
		suffix := ""
		if len(m) > 2 {
			suffix = "..."
		}
		return fmt.Sprintf("[%s%s](%s)", m[1], suffix, url)
	}
	return url
}
