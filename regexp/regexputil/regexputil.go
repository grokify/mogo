package regexputil

import (
	"regexp"
)

type RegexpSet struct {
	Regexps map[string]*regexp.Regexp
}

func NewRegexpSet() RegexpSet {
	set := RegexpSet{
		Regexps: map[string]*regexp.Regexp{}}
	return set
}

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
