package stringsutil

import (
	"strings"
)

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
