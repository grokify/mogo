package regexputil

import (
	"regexp"
)

// StringMatcher provides the ability to match a string
// against multiple regular expressions with a pre-formatting
// function as well.
type StringMatcher struct {
	Matchers []*regexp.Regexp
	PreMatch func(string) string
}

// NewStringMatcher returns a new StringMatcher instance.
func NewStringMatcher() StringMatcher {
	return StringMatcher{Matchers: []*regexp.Regexp{}}
}

// AddMatcher adds a regular expression to attempt to
// match against.
func (sm *StringMatcher) AddMatcher(rx *regexp.Regexp) {
	if sm.Matchers == nil {
		sm.Matchers = []*regexp.Regexp{}
	}
	sm.Matchers = append(sm.Matchers, rx)
}

// Match runs the provided string againts the prematch
// function and regular expresssions, returning true if
// any of the expressions match.
func (sm *StringMatcher) Match(input string) bool {
	if sm.PreMatch != nil {
		input = sm.PreMatch(input)
	}
	for _, rx := range sm.Matchers {
		if rx.MatchString(input) {
			return true
		}
	}
	return false
}
