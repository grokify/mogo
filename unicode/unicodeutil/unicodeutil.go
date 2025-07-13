package unicodeutil

import (
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// RemoveDiacritics removes diacritical marks (accents) from a Unicode string.
// A rune map can be used for manual overrides such as converting '™' to "tm".
func RemoveDiacritics(s string, m map[rune][]rune) string {
	t := norm.NFD.String(s)
	b := make([]rune, 0, len(t))

	for _, r := range t {
		if v, ok := m[r]; ok && len(v) > 0 {
			b = append(b, v...)
		} else if ok {
			continue
		} else if unicode.Is(unicode.Mn, r) {
			continue // skip mark (Mn) runes
		} else {
			b = append(b, r)
		}
	}

	return string(b)
}

// Unescape wraps the `strconv.Unquote()` function to provide a method of converting
// \u escaped unicode literals. It can take a string like "M\\u00fcnchen" and return "München"
func Unescape(input string) (string, error) {
	// Ensure proper escaping (e.g., turn \uXXXX into \\uXXXX if needed)
	escaped := input

	// Wrap in double quotes for strconv to parse correctly
	quoted := `"` + strings.ReplaceAll(escaped, `"`, `\"`) + `"`
	return strconv.Unquote(quoted)
}
