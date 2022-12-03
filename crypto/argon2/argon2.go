package argon2

import (
	"encoding/base32"
	"strings"

	"github.com/martinlindhe/base36"
	argon2lib "golang.org/x/crypto/argon2"
)

func HashSimple(input, salt []byte) []byte {
	return argon2lib.IDKey(input, salt, 1, 64*1024, 4, 32)
}

// HashSimple returns an argon2id hashed password encoded in Base36.
// The base36 library always returns upper case.
func HashSimpleBase36(input, salt []byte) string {
	return base36.EncodeBytes(HashSimple(input, salt))
}

// HashSimple returns an argon2id hashed password encoded in Base36.
// The base36 library always returns upper case.
func HashSimpleBase32(input, salt []byte, trim bool) string {
	if trim {
		return strings.Trim(base32.StdEncoding.EncodeToString(HashSimple(input, salt)), "=")
	}
	return base32.StdEncoding.EncodeToString(HashSimple(input, salt))
}
