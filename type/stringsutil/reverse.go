package stringsutil

import (
	"strings"
)

// Reverse reverses string using strings.Builder. It's about 3 times faster
// than the one with using a string concatenation
func Reverse(s string) string {
	var sb strings.Builder
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		sb.WriteRune(runes[i])
	}
	return sb.String()
}

// ReverseIndex returns the `Index` after reversing the
// supplied string and substring.
func ReverseIndex(s, substr string) int {
	return strings.Index(Reverse(s), Reverse(substr))
}
