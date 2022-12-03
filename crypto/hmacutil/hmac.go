// hmacutil provides HMAC utility functions.
package hmacutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
)

func HMACSHA256(key, msg []byte) []byte {
	sig := hmac.New(sha256.New, key)
	sig.Write(msg)
	return sig.Sum(nil)
}

func HMACSHA256Base32(key, msg []byte) string {
	// return bigint.MustEncodeToString(32, HMACSHA256(key, msg))
	return base32.StdEncoding.EncodeToString(HMACSHA256(key, msg))
}

func HMACSHA256Hex(key, msg []byte) string {
	return hex.EncodeToString(HMACSHA256(key, msg))
}

// Validate compares MACs for equality without leaking timing information.
func Validate(key, msg, msgMAC []byte) bool {
	return hmac.Equal(msgMAC, HMACSHA256(key, msg))
}

// ValidateBase32 compares MACs for equality without leaking timing information.
func ValidateBase32(key, msg []byte, msgMAC string) (bool, error) {
	msgMacBytes, err := base32.StdEncoding.DecodeString(msgMAC)
	if err != nil {
		return false, err
	}
	return hmac.Equal(msgMacBytes, HMACSHA256(key, msg)), nil
}

// ValidateHex compares MACs for equality without leaking timing information.
func ValidateHex(key, msg []byte, msgMAC string) (bool, error) {
	msgMacBytes, err := hex.DecodeString(msgMAC)
	if err != nil {
		return false, err
	}
	return hmac.Equal(msgMacBytes, HMACSHA256(key, msg)), nil
}
