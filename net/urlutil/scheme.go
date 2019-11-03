package urlutil

import (
	"fmt"
	"regexp"
	"strings"
)

const uriSchemePattern string = `^(?i)([a-z][0-9a-z\-\+.]+)://`

var rxScheme *regexp.Regexp = regexp.MustCompile(uriSchemePattern)

// https://www.iana.org/assignments/uri-schemes/uri-schemes.xhtml

// UriHasScheme returns a boolean true or false if the string
// has a URI scheme.
func UriHasScheme(uri string) bool {
	scheme := UriScheme(uri)
	if len(scheme) > 0 {
		return true
	}
	return false
}

// UriScheme extracts the URI scheme from a string. It returns
// an empty string if none is encountered.
func UriScheme(uri string) string {
	uri = strings.TrimSpace(uri)
	m := rxScheme.FindAllStringSubmatch(uri, -1)
	if len(m) > 0 && len(m[0]) == 2 {
		fmt.Println(m[0][1])
		return strings.TrimSpace(m[0][1])
	}
	return ""
}
