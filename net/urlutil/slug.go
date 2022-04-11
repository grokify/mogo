package urlutil

import (
	"regexp"
)

var (
	rxSlugTextToURL *regexp.Regexp = regexp.MustCompile(`[\s\t\r\n_-]+`)
	rxSlugURLToText *regexp.Regexp = regexp.MustCompile(`[_-]+`)
)

func SlugTextToURL(s string, underscore bool) string {
	sep := "-"
	if underscore {
		sep = "_'"
	}
	return rxSlugTextToURL.ReplaceAllString(s, sep)
}

func SlugURLToText(s string) string {
	return rxSlugURLToText.ReplaceAllString(s, " ")
}
