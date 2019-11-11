package urlutil

import (
	"net/url"
	"regexp"
)

const (
	netURLTemplateSubstituteHostname string = "Iaae8QKUZsUZpxDjpkV0GII9ahZLTeWiyFasAO8xPmv"
	urlPattern                       string = `^([0-9a-zA-Z]+://)([^/]+)`
)

var urlPatternRx *regexp.Regexp = regexp.MustCompile(urlPattern)

// ParseURLTemplate exists to parse templates with variables
// that do not meet RFC specifications. For example:
// https://{customer}.example.com:{port}/v5
// "invalid URL escape "%7B"" for `{` within a Hostname or
// "invalid port ":{port}" after host"
func ParseURLTemplate(input string) (*url.URL, error) {
	workingURL := input
	m := urlPatternRx.FindStringSubmatch(workingURL)
	if len(m) == 0 {
		return url.Parse(workingURL)
	}

	replString := m[1] + netURLTemplateSubstituteHostname
	workingURL = urlPatternRx.ReplaceAllString(workingURL, replString)

	parsedURL, err := url.Parse(workingURL)
	if err != nil {
		return nil, err
	}
	parsedURL.Host = m[2]

	return parsedURL, nil
}
