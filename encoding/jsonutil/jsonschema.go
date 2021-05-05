package jsonutil

import (
	"fmt"
	"regexp"
)

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

var rxSlashMore = regexp.MustCompile(`/+`)

// PointerCondense removes duplicate slashes.
func PointerCondense(s string) string {
	return rxSlashMore.ReplaceAllString(s, "/")
}

// PointerSubEscapeAll will substitute vars using `fmt.Sprintf()`
// All strings are escaped.
func PointerSubEscapeAll(format string, vars ...interface{}) string {
	if len(vars) == 0 {
		return format
	}
	for i, v := range vars {
		if vString, ok := v.(string); ok {
			vars[i] = PropertyNameEscape(vString)
		}
	}
	return PointerCondense(fmt.Sprintf(format, vars...))
}
