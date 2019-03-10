package stringsutil

import (
	"regexp"
	"strings"
)

type MatchType int

const (
	Exact MatchType = iota
	TrimSpace
	TrimSpaceLower
	Regexp
)

type MatchInfo struct {
	MatchType MatchType
	String    string
	Regexp    *regexp.Regexp
}

// Match provides an canonical way to match strings using multiple
// approaches
func Match(s string, matchInfo MatchInfo) bool {
	switch matchInfo.MatchType {
	case Exact:
		if s == matchInfo.String {
			return true
		}
		return false
	case TrimSpace:
		m := strings.TrimSpace(s)
		if m == strings.TrimSpace(matchInfo.String) {
			return true
		}
		return false
	case TrimSpaceLower:
		m := strings.ToLower(strings.TrimSpace(s))
		if m == strings.ToLower(strings.TrimSpace(matchInfo.String)) {
			return true
		}
		return false
	case Regexp:
		if matchInfo.Regexp == nil {
			return false
		}
		return matchInfo.Regexp.MatchString(s)
	default:
		return false
	}
}
