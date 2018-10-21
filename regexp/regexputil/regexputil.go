package regexputil

import (
	"regexp"
)

// McReplaceAllString is a single line MustCompile regexp for ReplaceAllString

func McReplaceAllString(pattern string, s string, repl string) string {
	return regexp.MustCompile(pattern).ReplaceAllString(s, repl)
}

// RegexpSet is a struct that holds compiled regular expressions.
// Primary goals of this struct is to reduce MustCompile regular
// expressions into a single function call and to store the compiled
// regular expressions if desired

type RegexpSet struct {
	Regexps map[string]*regexp.Regexp
}

// NewRegexpSet returns a new RegexpSet struct

func NewRegexpSet() RegexpSet {
	set := RegexpSet{
		Regexps: map[string]*regexp.Regexp{}}
	return set
}

func (set *RegexpSet) GetRegexp(pattern string, useStore bool, key string) *regexp.Regexp {
	var rx *regexp.Regexp
	if useStore {
		if len(key) == 0 {
			key = pattern
		}
		var ok bool
		rx, ok = set.Regexps[key]
		if !ok {
			rx = regexp.MustCompile(pattern)
		}
		set.Regexps[key] = rx
	} else {
		rx = regexp.MustCompile(pattern)
	}
	return rx
}

func (set *RegexpSet) FindAllString(pattern string, s string, n int, useStore bool, key string) []string {
	rx := set.GetRegexp(pattern, useStore, key)
	rs := rx.FindAllString(s, n)
	return rs
}

// FindAllStringSubmatch performs a regular expression find against the
// supplied pattern and string. It will store the compiled regular expression
// for later use.

func (set *RegexpSet) FindAllStringSubmatch(pattern string, s string, n int, useStore bool, key string) [][]string {
	rx := set.GetRegexp(pattern, useStore, key)
	rs := rx.FindAllStringSubmatch(s, n)
	return rs
}

func (set *RegexpSet) FindStringSubmatch(pattern string, s string, useStore bool, key string) []string {
	rx := set.GetRegexp(pattern, useStore, key)
	rs := rx.FindStringSubmatch(s)
	return rs
}

func FindStringSubmatchNamedMap(rx *regexp.Regexp, s string) map[string]string {
	match := rx.FindStringSubmatch(s)
	result := make(map[string]string)
	for i, name := range rx.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result
}
