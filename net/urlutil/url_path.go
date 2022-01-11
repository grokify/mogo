package urlutil

import (
	"errors"
	"net/url"
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

func GetPathLeaf(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	sep := "/"
	p := strings.Trim(u.Path, sep)
	parts := strings.Split(p, sep)
	if len(parts) == 0 {
		return "", errors.New("GetPathLeaf - no path")
	}
	return parts[len(parts)-1], nil
}

func ModifyPath(rawurl, newpath string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	newpath = strings.TrimSpace(newpath)
	if newpath == "/" {
		newpath = ""
	}
	u.Path = newpath
	return CondenseUri(u.String()), nil
}
