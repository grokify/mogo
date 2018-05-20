package argon2

import (
	"github.com/martinlindhe/base36"
	argon2lib "golang.org/x/crypto/argon2"
)

// HashSimple returns an argon2id hashed password encoded in Base36.
// The base36 library always returns upper case.
func HashSimpleBase36(input, salt []byte) string {
	return base36.EncodeBytes(argon2lib.IDKey(input, salt, 1, 64*1024, 4, 32))
}
