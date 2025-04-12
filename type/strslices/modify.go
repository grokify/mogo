package strslices

import "strings"

func Map(s []string, fn func(s string) string) []string {
	var out []string
	if fn == nil {
		fn = func(s string) string { return s }
	}
	for _, si := range s {
		out = append(out, fn(si))
	}
	return out
}

func ToLower(s []string) []string {
	return Map(s, strings.ToLower)
}

func ToUpper(s []string) []string {
	return Map(s, strings.ToUpper)
}
