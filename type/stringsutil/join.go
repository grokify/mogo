package stringsutil

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func JoinTrimSpace(strs []string) string {
	return rxSpaces.ReplaceAllString(strings.Join(strs, " "), " ")
}

// JoinInterface joins an interface and returns a string. It takes
// a join separator, boolean to replace the join separator in the
// string parts and a separator alternate. `stripEmbeddedSep` strips
// separator string found within parts. `stripRepeatedSep` strips
// repeating separators. This flexibility is designed to support
// joining data for both CSVs and paths.
func JoinInterface(arr []any, sep string, stripRepeatedSep bool, stripEmbeddedSep bool, altSep string) string {
	parts := []string{}
	rx := regexp.MustCompile(sep)
	for _, el := range arr {
		part := fmt.Sprintf("%v", el)
		if stripEmbeddedSep {
			part = rx.ReplaceAllString(part, altSep)
		}
		parts = append(parts, part)
	}
	joined := strings.Join(parts, sep)
	if stripRepeatedSep {
		joined = regexp.MustCompile(fmt.Sprintf("%s+", sep)).
			ReplaceAllString(joined, sep)
	}
	return joined
}

func JoinLiterary(slice []string, sep, joinWord string) string {
	switch len(slice) {
	case 0:
		return ""
	case 1:
		return slice[0]
	case 2:
		return slice[0] + " " + joinWord + " " + slice[1]
	default:
		last, rest := slice[len(slice)-1], slice[:len(slice)-1]
		rest = append(rest, joinWord+" "+last)
		return strings.Join(rest, sep+" ")
	}
}

func JoinLiteraryQuote(slice []string, leftQuote, rightQuote, sep, joinWord string) string {
	newSlice := SliceCondenseAndQuoteSpace(slice, leftQuote, rightQuote)
	switch len(newSlice) {
	case 0:
		return ""
	case 1:
		return newSlice[0]
	case 2:
		return newSlice[0] + " " + joinWord + " " + newSlice[1]
	default:
		last, rest := newSlice[len(newSlice)-1], newSlice[:len(newSlice)-1]
		rest = append(rest, joinWord+" "+last)
		return strings.Join(rest, sep+" ")
	}
}

func JoinStringsTrimSpaceToLowerSort(strs []string, sep string) string {
	wip := []string{}
	for _, s := range strs {
		s = strings.ToLower(strings.TrimSpace(s))
		if len(s) > 0 {
			wip = append(wip, s)
		}
	}
	sort.Strings(wip)
	return strings.Join(wip, sep)
}
