package stringsutil

import "strings"

func SuffixStrip(s, suffix string) string {
	if ReverseIndex(s, suffix) == 0 {
		return s[:len(s)-len(suffix)]
	}
	return s
}

func SuffixReplace(s, oldSuffix, newSuffix string) string {
	if len(oldSuffix) == 0 {
		return s + newSuffix
	} else if ReverseIndex(s, oldSuffix) == 0 {
		return SuffixStrip(s, oldSuffix) + newSuffix
	}
	return s
}

func SuffixParse(s, wantSuffix string) (fullstring, prefix, suffix string) {
	fullstring = s
	if len(suffix) > len(fullstring) {
		return
	}
	lastIndex := strings.LastIndex(s, wantSuffix)
	if lastIndex > 0 && lastIndex == len(s)-len(wantSuffix) {
		prefix = fullstring[:len(s)-len(wantSuffix)]
		suffix = wantSuffix
		return
	}
	return
}

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
