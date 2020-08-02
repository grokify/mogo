package stringsutil

import (
	"regexp"
)

var rxNonDigit = regexp.MustCompile(`[^0-9]`)

func DigitsOnly(input string) string {
	return rxNonDigit.ReplaceAllString(input, "")
}
