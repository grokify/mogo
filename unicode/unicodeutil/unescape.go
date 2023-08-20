package unicodeutil

import (
	"fmt"
	"strconv"
	"strings"
)

// Unescape wraps the `strconv.Unquote()` function to provide a method of converting
// \u escaped unicode literals.
func Unescape(s string) (string, error) {
	if !strings.Contains(s, QuotationMark) {
		return strconv.Unquote(QuotationMark + s + QuotationMark)
	} else if !strings.Contains(s, Apostrophe) {
		return strconv.Unquote(Apostrophe + s + Apostrophe)
	} else {
		return s, fmt.Errorf("cannot unescape string (%s)", s)
	}
}
