package regexputil

import (
	"regexp"
)

// RegexpSet is a struct that holds compiled regular expressions
type RegexpSet struct {
	Regexps map[string]*regexp.Regexp
}

// NewRegexpSet returns a new RegexpSet struct
func NewRegexpSet() RegexpSet {
	set := RegexpSet{
		Regexps: map[string]*regexp.Regexp{}}
	return set
}

// FindAllStringSubmatch performs a regular expression find against the
// supplied pattern and string. It will store the compiled regular expression
// for later use.
func (set *RegexpSet) FindAllStringSubmatch(pattern string, s string, key string) [][]string {
	if len(key) == 0 {
		key = pattern
	}
	rx, ok := set.Regexps[key]
	if !ok {
		rx = regexp.MustCompile(pattern)
	}
	set.Regexps[key] = rx
	rs := rx.FindAllStringSubmatch(s, -1)
	return rs
}
