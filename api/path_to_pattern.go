package api

import (
	"regexp"
	"strings"
)

const rxMatchParameterPattern = `{[^\{\}]*}`
const rxMatchParameterActual = `[^\{\}]*`

// URLTransformer is useful for reading log files and converting actual
// request URls into pattners, such as those used in the OpenAPI Spec.
type URLTransformer struct {
	ExactPaths     []string
	RegexpPaths    map[string]*regexp.Regexp
	rxMatchPattern *regexp.Regexp
	rxMatchActual  *regexp.Regexp
	rxStripQuery   *regexp.Regexp
}

func NewURLTransformer() URLTransformer {
	return URLTransformer{
		ExactPaths:     []string{},
		RegexpPaths:    map[string]*regexp.Regexp{},
		rxMatchPattern: regexp.MustCompile(rxMatchParameterPattern),
		rxMatchActual:  regexp.MustCompile(rxMatchParameterActual),
		rxStripQuery:   regexp.MustCompile(`\?.*$`)}
}

// LoadPaths to lost config data. See the test file for an example.
func (t *URLTransformer) LoadPaths(paths []string) error {
	for _, path := range paths {
		err := t.LoadPath(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *URLTransformer) LoadPath(path string) error {
	path = t.rxStripQuery.ReplaceAllString(path, "")
	i1 := strings.Index(path, "{")
	i2 := strings.Index(path, "}")
	if i1 < 0 && i2 < 0 {
		t.ExactPaths = append(t.ExactPaths, path)
		return nil
	}
	linkPattern := t.rxMatchPattern.ReplaceAllString(path, rxMatchParameterActual)
	linkPattern = `^` + linkPattern + `$`
	rx, err := regexp.Compile(linkPattern)
	if err != nil {
		return err
	}
	t.RegexpPaths[path] = rx
	return nil
}

// URLActualToPattern is the "runtime" API that is called over and over
// for RUL classification purposes.
func (t *URLTransformer) URLActualToPattern(s string) string {
	s = t.rxStripQuery.ReplaceAllString(s, "")
	for _, try := range t.ExactPaths {
		if s == try {
			return s
		}
	}
	for pattern, rx := range t.RegexpPaths {
		if rx.MatchString(s) {
			return pattern
		}
	}
	return s
}
