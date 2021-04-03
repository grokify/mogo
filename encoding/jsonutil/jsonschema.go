package jsonutil

import "regexp"

const (
	EscapedTilde = "~0"
	EscapedSlash = "~1"
)

var rxTilde = regexp.MustCompile(`~`)
var rxSlash = regexp.MustCompile(`/`)
var rxTildeEscaped = regexp.MustCompile(`~0`)
var rxSlashEscaped = regexp.MustCompile(`~1`)

// PropertyNameEscape escapes JSON property name using
// JSON Schema rules.
func PropertyNameEscape(s string) string {
	// https://opis.io/json-schema/1.x/pointers.html
	// https://json-schema.org/understanding-json-schema/structuring.html
	return rxSlash.ReplaceAllString(
		rxTilde.ReplaceAllString(s, EscapedTilde),
		EscapedSlash)
}

// PropertyNameUnescape unescapes JSON property name using
// JSON Schema rules.
func PropertyNameUnescape(s string) string {
	// https://opis.io/json-schema/1.x/pointers.html
	// https://json-schema.org/understanding-json-schema/structuring.html
	return rxTildeEscaped.ReplaceAllString(
		rxSlashEscaped.ReplaceAllString(s, "/"),
		"~")
}
