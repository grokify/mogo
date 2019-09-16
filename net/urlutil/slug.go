package urlutil

import (
	"regexp"
)

var (
	rxSlugTextToUrl *regexp.Regexp = regexp.MustCompile(`[\s]+`)
	rxSlugUrlToText *regexp.Regexp = regexp.MustCompile(`[_]+`)
)

func SlugTextToUrl(s string) string {
	return rxSlugTextToUrl.ReplaceAllString(s, "_")
}

func SlugUrlToText(s string) string {
	return rxSlugUrlToText.ReplaceAllString(s, " ")
}
