package stringsutil

import (
	"regexp"
	"strings"
	"time"
)

type MatchType int

const (
	MatchExact MatchType = iota
	MatchTrimSpace
	MatchTrimSpaceLower
	MatchRegexp
	MatchTimeGT
	MatchTimeGTE
	MatchTimeLT
	MatchTimeLTE
	MatchHasSuffix
	MatchHasPrefix
)

type MatchInfo struct {
	MatchType  MatchType
	String     string
	Regexp     *regexp.Regexp
	TimeLayout string
	TimeMin    time.Time
	TimeMax    time.Time
}

// Match provides an canonical way to match strings using multiple
// approaches
func Match(s string, matchInfo MatchInfo) (bool, error) {
	switch matchInfo.MatchType {
	case MatchExact:
		if s == matchInfo.String {
			return true, nil
		}
		return false, nil
	case MatchTrimSpace:
		m := strings.TrimSpace(s)
		if m == strings.TrimSpace(matchInfo.String) {
			return true, nil
		}
		return false, nil
	case MatchTrimSpaceLower:
		m := strings.ToLower(strings.TrimSpace(s))
		if m == strings.ToLower(strings.TrimSpace(matchInfo.String)) {
			return true, nil
		}
		return false, nil
	case MatchRegexp:
		if matchInfo.Regexp == nil {
			return false, nil
		}
		return matchInfo.Regexp.MatchString(s), nil
	case MatchTimeGT:
		t, err := time.Parse(matchInfo.TimeLayout, s)
		if err != nil {
			return false, err
		}
		if t.After(matchInfo.TimeMin) {
			return true, nil
		} else {
			return false, nil
		}
	case MatchTimeGTE:
		t, err := time.Parse(matchInfo.TimeLayout, s)
		if err != nil {
			return false, err
		}
		if t.After(matchInfo.TimeMin) || t.Equal(matchInfo.TimeMin) {
			return true, nil
		} else {
			return false, nil
		}
	case MatchTimeLT:
		t, err := time.Parse(matchInfo.TimeLayout, s)
		if err != nil {
			return false, err
		}
		if t.Before(matchInfo.TimeMax) {
			return true, nil
		} else {
			return false, nil
		}
	case MatchTimeLTE:
		t, err := time.Parse(matchInfo.TimeLayout, s)
		if err != nil {
			return false, err
		}
		if t.Before(matchInfo.TimeMax) || t.Equal(matchInfo.TimeMax) {
			return true, nil
		} else {
			return false, nil
		}
	default:
		return false, nil
	}
}
