package boolutil

import (
	"strings"
)

func StringToBool(value string) bool {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	if value == "true" {
		return true
	}
	return false
}
