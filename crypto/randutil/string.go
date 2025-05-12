package randutil

import "github.com/grokify/mogo/encoding/basex"

// RandString returns a random string of length `length` using the supplied alphabet.
// If no alphabet is provided, `AlphabetBase16`, aha hexadecimal is used.
func RandString(alphabet string, length uint) (string, error) {
	if length == 0 {
		return "", nil
	}
	if len(alphabet) == 0 {
		alphabet = basex.AlphabetBase16
	}
	var out string
	for i := uint(0); i < length; i++ {
		idx, err := CryptoRandIntInRange(0, len(alphabet)-1)
		if err != nil {
			return "", err
		} else {
			out += string(alphabet[idx])
		}
	}
	return out, nil
}
