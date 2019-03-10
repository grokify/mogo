package stringsutil

import (
	"strings"
)

// StringToBool converts a string to a boolean value
// looking for the string "true" in any case.
func ToBool(v string) bool {
	if strings.TrimSpace(strings.ToLower(v)) == "true" {
		return true
	}
	return false
}
