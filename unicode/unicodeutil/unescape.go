package unicodeutil

import (
	"strconv"
	"strings"
)

// Unescape wraps the `strconv.Unquote()` function to provide a method of converting
// \u escaped unicode literals. It can take a string like "M\\u00fcnchen" and return "München"
func Unescape(input string) (string, error) {
	// Ensure proper escaping (e.g., turn \uXXXX into \\uXXXX if needed)
	escaped := input

	// Wrap in double quotes for strconv to parse correctly
	quoted := `"` + strings.ReplaceAll(escaped, `"`, `\"`) + `"`
	return strconv.Unquote(quoted)
}
