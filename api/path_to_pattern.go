package api

import (
	"regexp"
	"strings"

	"github.com/grokify/gotilla/type/stringsutil/join"
)

const (
	rxMatchParameterPattern = `{[^\{\}/]*?}`
	rxMatchParameterActual  = `[^\{\}/]*?`
)

// URLTransformer is useful for reading log files and converting actual
// request URls into pattners, such as those used in the OpenAPI Spec for
// reporting and categorization purposes.
type URLTransformer struct {
	ExactPaths     []string
	RegexpPaths    map[string]*regexp.Regexp
	rxMatchPattern *regexp.Regexp
	rxMatchActual  *regexp.Regexp
	rxStripQuery   *regexp.Regexp
}

// NewURLTransformer creates a new URLTransformer instance.
func NewURLTransformer() URLTransformer {
	return URLTransformer{
		ExactPaths:     []string{},
		RegexpPaths:    map[string]*regexp.Regexp{},
		rxMatchPattern: regexp.MustCompile(rxMatchParameterPattern),
		rxMatchActual:  regexp.MustCompile(rxMatchParameterActual),
		rxStripQuery:   regexp.MustCompile(`\?.*$`)}
}

// LoadPaths loads multiple spec URL patterns. See the test file for an example.
func (ut *URLTransformer) LoadPaths(paths []string) error {
	for _, path := range paths {
		err := ut.LoadPath(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadPath loads a single spec URL pattern.
func (ut *URLTransformer) LoadPath(path string) error {
	path = ut.rxStripQuery.ReplaceAllString(path, "")
	i1 := strings.Index(path, "{")
	i2 := strings.Index(path, "}")
	if i1 < 0 && i2 < 0 {
		ut.ExactPaths = append(ut.ExactPaths, path)
		return nil
	}
	linkPattern := ut.rxMatchPattern.ReplaceAllString(path, rxMatchParameterActual)
	linkPattern = `^` + linkPattern + `$`
	rx, err := regexp.Compile(linkPattern)
	if err != nil {
		return err
	}
	ut.RegexpPaths[path] = rx
	return nil
}

// URLActualToPattern is the "runtime" API that is called over and over
// for URL classification purposes.
func (ut *URLTransformer) URLActualToPattern(s string) string {
	s = ut.rxStripQuery.ReplaceAllString(s, "")
	for _, try := range ut.ExactPaths {
		if s == try {
			return s
		}
	}
	for pattern, rx := range ut.RegexpPaths {
		if rx.MatchString(s) {
			return pattern
		}
	}
	return s
}

func (ut *URLTransformer) BuildReverseEndpointPattern(method, actualURL string) string {
	pattern := ut.URLActualToPattern(actualURL)
	return join.JoinCondenseTrimSpace([]string{pattern, method}, " ")
}
