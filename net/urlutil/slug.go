package urlutil

import (
	"regexp"
)

var (
	rxSlugTextToURL *regexp.Regexp = regexp.MustCompile(`[\s\t\r\n]+`)
	rxSlugURLToText *regexp.Regexp = regexp.MustCompile(`[_]+`)
)

func SlugTextToURL(s string) string {
	return rxSlugTextToURL.ReplaceAllString(s, "_")
}

func SlugURLToText(s string) string {
	return rxSlugURLToText.ReplaceAllString(s, " ")
}
