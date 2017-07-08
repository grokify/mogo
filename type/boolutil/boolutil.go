package boolutil

import (
	"strings"
)

// StringToBool converts a string to a boolean value
// looking for the string "true" in any case.
func StringToBool(v string) bool {
	if strings.TrimSpace(strings.ToLower(v)) == "true" {
		return true
	}
	return false
}
