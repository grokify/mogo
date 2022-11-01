package stringsutil

import (
	"regexp"
)

var rxNonDigit = regexp.MustCompile(`[^0-9]`)

func DigitsOnly(input string) string {
	return rxNonDigit.ReplaceAllString(input, "")
}

var rxDigits = regexp.MustCompile(`^[0-9]+$`)

func IsInteger(input string) bool {
	return rxDigits.MatchString(input)
}
