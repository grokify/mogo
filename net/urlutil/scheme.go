package urlutil

import (
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
		return strings.TrimSpace(m[0][1])
	}
	return ""
}

func IsHttp(uri string, inclHttp, inclHttps bool) bool {
	try := strings.ToLower(strings.TrimSpace(uri))
	if strings.Index(try, "http://") == 0 && inclHttp {
		return true
	} else if strings.Index(try, "https://") == 0 && inclHttps {
		return true
	}
	return false
}
