package urlutil

import (
	"regexp"
	"strings"
)

var leading *regexp.Regexp = regexp.MustCompile(`^/+`)
var trailing *regexp.Regexp = regexp.MustCompile(`/+$`)

// SplitPath splits a URL path string with optional removal of leading
// and trailing slashes.
func SplitPath(urlPath string, stripLeading, stripTrailing bool) []string {
	urlPath = strings.TrimSpace(urlPath)
	if stripLeading {
		urlPath = leading.ReplaceAllString(urlPath, "")
	}
	if stripTrailing {
		urlPath = trailing.ReplaceAllString(urlPath, "")
	}
	return strings.Split(urlPath, "/")
}
