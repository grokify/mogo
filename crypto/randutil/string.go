package randutil

import "errors"

const (
	AlphabetBase10 = "0123456789"
	AlphabetBase16 = "0123456789abcdef"
	AlphabetBase36 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphabetBase62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func RandString(alphabet string, length uint) (string, error) {
	if length == 0 {
		return "", nil
	}
	if len(alphabet) == 0 {
		alphabet = AlphabetBase16
	}
	rnd := NewCryptoRand(nil, nil)
	out := ""
	for i := 0; i < int(length); i++ {
		idx, err := rnd.Intn(len(alphabet))
		if err != nil {
			return "", err
		}
		if idx >= len(alphabet) {
			return "", errors.New("idx too large")
		}
		out += string(alphabet[idx])
	}
	return out, nil
}

func MustRandString(alphabet string, length uint) string {
	if length == 0 {
		return ""
	}
	if len(alphabet) == 0 {
		alphabet = AlphabetBase16
	}
	rnd := NewCryptoRand(nil, nil)
	out := ""
	for i := 0; i < int(length); i++ {
		idx := rnd.MustIntn(len(alphabet))
		if idx >= len(alphabet) {
			panic("idx too large")
		}
		out += string(alphabet[idx])
	}
	return out
}
