package stringsutil

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type MatchType int

const (
	MatchExact MatchType = iota
	MatchStringTrimSpace
	MatchStringTrimSpaceLower
	MatchStringPrefix
	MatchStringSuffix
	MatchStringRegexp
	MatchTimeGT
	MatchTimeGTE
	MatchTimeLT
	MatchTimeLTE
)

type MatchInfo struct {
	MatchType  MatchType
	String     string
	Regexp     *regexp.Regexp
	TimeLayout string
	TimeMin    time.Time
	TimeMax    time.Time
}

// Match provides an canonical way to match strings using multiple approaches
func Match(s string, matchInfo MatchInfo) (bool, error) {
	switch matchInfo.MatchType {
	case MatchExact:
		return s == matchInfo.String, nil
	case MatchStringTrimSpace:
		return strings.TrimSpace(s) == strings.TrimSpace(matchInfo.String), nil
	case MatchStringTrimSpaceLower:
		return strings.EqualFold(strings.TrimSpace(s), strings.TrimSpace(matchInfo.String)), nil
	case MatchStringPrefix:
		return strings.Index(s, matchInfo.String) == 0, nil
	case MatchStringSuffix:
		return ReverseIndex(s, matchInfo.String) == 0, nil
	case MatchStringRegexp:
		if matchInfo.Regexp == nil {
			return false, errors.New("no regexp supplied")
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

/*
func CheckSuffix(input, wantSuffix string) (fullstring, prefix, suffix string) {
	fullstring = input
	if len(suffix) > len(fullstring) {
		return
	}
	lastIndex := strings.LastIndex(input, wantSuffix)
	if lastIndex > 0 && lastIndex == len(input)-len(wantSuffix) {
		prefix = fullstring[:len(input)-len(wantSuffix)]
		suffix = wantSuffix
		return
	}
	return
}*/

/*
func SuffixMap(inputs, suffixes []string) (prefixes []string, matches map[string]string, nonmatches []string) {
	matches = map[string]string{}
	for _, input := range inputs {
		gotMatch := false
		for _, suffix := range suffixes {
			gotFull, gotPrefix, gotSuffix := SuffixParse(input, suffix)
			if suffix == gotSuffix {
				matches[gotSuffix] = gotFull
				prefixes = append(prefixes, gotPrefix)
				gotMatch = true
			}
		}
		if !gotMatch {
			nonmatches = append(nonmatches, input)
		}
	}
	prefixes = SliceCondenseSpace(prefixes, true, true)
	nonmatches = SliceCondenseSpace(nonmatches, true, true)
	return prefixes, matches, nonmatches
}
*/
