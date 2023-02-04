package randutil

import "github.com/grokify/mogo/encoding"

// RandString returns a random string of length `length` using the supplied alphabet.
// If no alphabet is provided, `AlphabetBase16`, aha hexadecimal is used.
func RandString(alphabet string, length uint) string {
	if length == 0 {
		return ""
	}
	if len(alphabet) == 0 {
		alphabet = encoding.AlphabetBase16
	}
	var out string
	for i := 0; i < int(length); i++ {
		idx := Intn(uint(len(alphabet)))
		if idx >= len(alphabet) {
			panic("idx too large")
		}
		out += string(alphabet[idx])
	}
	return out
}
