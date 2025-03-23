package stringsutil

import (
	"regexp"
	"strings"
)

var rxSplitLines = regexp.MustCompile(`(\r\n|\r|\n)`)

// SplitLines splits a string on the regxp `(\r\n|\r|\n)`.
func SplitLines(text string) []string {
	return rxSplitLines.Split(text, -1)
}

// SplitTrimSpace splits a string and trims spaces on
// remaining elements, and optionally removing empty elements.
func SplitTrimSpace(s, sep string, exclEmptyString bool) []string {
	strs := []string{}
	split := strings.Split(s, sep)
	for _, str := range split {
		str = strings.TrimSpace(str)
		if !exclEmptyString || str != "" {
			strs = append(strs, str)
		}
	}
	return strs
}

// SplitCSVUnescaped splits a string on commas without a following space. It is useful
// when where commas are not escaped. An example is joining two strings `["foo, bar", "baz, qux"]`
// as a CSV such as `foo, bar,baz, qux`.
func SplitCSVUnescaped(s string) []string {
	s = strings.ReplaceAll(s, ", ", "%2C ")
	p := strings.Split(s, ",")
	for i, px := range p {
		p[i] = strings.ReplaceAll(px, "%2C ", ", ")
	}
	return p
}

/*
// EOL: SplitLines has been removed in favor of `SplitLines()`.
func SplitLines(s string) []string {
	return strings.Split(ToLineFeeds(s), "\n")
}
*/

/*
// SplitTrimSpace splits a string and trims spaces on remaining elements.
// EOL: SplitTrimSpace has been removed in favor of renaming
// `SplitCondenseSpace()` to `SplitTrimSpace()` with `exclEmptyString` param.
func SplitTrimSpace(s, sep string) []string {
	split := strings.Split(s, sep)
	strs := []string{}
	for _, str := range split {
		strs = append(strs, strings.TrimSpace(str))
	}
	return strs
}
*/
