package urlutil

import (
	"regexp"
	"strings"
)

var rxPathVarToGeneric = regexp.MustCompile(`{[^}{]*}`)

func VarsToGeneric(input string) string {
	return rxPathVarToGeneric.ReplaceAllString(input, "{}")
}

func MatchGeneric(path1, path2 string) bool {
	gen1 := VarsToGeneric(path1)
	gen2 := VarsToGeneric(path2)
	if gen1 != gen2 {
		return false
	}
	return true
}

func EndpointString(path, method string, generic bool) string {
	path = strings.TrimSpace(path)
	method = strings.ToUpper(strings.TrimSpace(method))
	parts := []string{}
	if len(path) > 0 {
		if generic {
			path = VarsToGeneric(path)
		}
		parts = append(parts, path)
	}
	if len(method) > 0 {
		parts = append(parts, method)
	}
	return strings.Join(parts, " ")
}
