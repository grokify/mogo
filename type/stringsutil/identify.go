package stringsutil

import (
	"regexp"
)

func DigitsOnly(input string) string {
	return rxNonNumeric.ReplaceAllString(input, "")
}

var (
	rxAlpha           = regexp.MustCompile(`^[A-Za-z]+$`)
	rxAlphaNumeric    = regexp.MustCompile(`^[0-9A-Za-z]+$`)
	rxAlphaNumericNot = regexp.MustCompile(`[^0-9A-Za-z]`)
	rxNumeric         = regexp.MustCompile(`^[0-9]+$`)
	rxNonNumeric      = regexp.MustCompile(`[^0-9]`)
)

func IsAlpha(s string) bool {
	return rxAlpha.MatchString(s)
}

func IsAlphaNumeric(s string) bool {
	return rxAlphaNumeric.MatchString(s)
}

func IsNumeric(s string) bool {
	return rxNumeric.MatchString(s)
}

func ToAlphaNumeric(s string) string {
	return rxAlphaNumericNot.ReplaceAllString(s, "")
}
