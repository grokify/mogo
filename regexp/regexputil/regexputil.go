package regexputil

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/grokify/mogo/strconv/strconvutil"
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
		name = strings.TrimSpace(name)
		if i != 0 && i < len(match) && name != "" {
			result[name] = match[i]
		}
	}
	return result
}

// CaptureUint requires first capture to be parseable by `strconv.Atoi` such as `[0-9]+`.
// Panics is capture is not parseable.
func CaptureUint(b []byte, expr string) (uint, error) {
	m := regexp.MustCompile(expr).FindSubmatch(b)
	if len(m) > 1 {
		vi, err := strconvutil.Atou(string(m[1]))
		if err != nil {
			return 0, err
		}
		return vi, nil
	}
	return 0, ErrMatchNotFound
}

// MustCaptureUint will return `0` if match is not found. It will panic if strconv fails.
func MustCaptureUint(b []byte, expr string) uint {
	v, err := CaptureUint(b, expr)
	if err != nil {
		panic(err)
	}
	return v
}

// CaptureString returns the first captured string.
func CaptureString(b []byte, expr string) string {
	rx := regexp.MustCompile(expr)
	m := rx.FindSubmatch(b)
	if len(m) > 1 {
		return string(m[1])
	}
	return ""
}

var ErrMatchNotFound = errors.New("match not found")

func CaptureTime(b []byte, expr, layout string) (time.Time, error) {
	rx := regexp.MustCompile(expr)
	m := rx.FindSubmatch(b)
	if len(m) > 1 {
		return time.Parse(layout, string(m[1]))
	}
	return time.Time{}, ErrMatchNotFound
}

func MustCaptureTime(b []byte, expr, layout string) time.Time {
	t, err := CaptureTime(b, expr, layout)
	if err != nil {
		panic(err)
	}
	return t
}
