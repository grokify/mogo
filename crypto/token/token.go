package token

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/base58"
)

// ---- public encoding helpers ----

// Base32 returns n random bytes encoded as Base32.
func Base32(n int) (string, error) {
	return fromEncoder(n, func(b []byte) string {
		return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	})
}

// Base58 returns n random bytes encoded in Bitcoin Base58.
func Base58(n int) (string, error) {
	return fromEncoder(n, func(b []byte) string {
		return base58.Encode(b)
	})
}

// Base64 returns n random bytes encoded as Base64.
func Base64(n int) (string, error) {
	return fromEncoder(n, func(b []byte) string {
		return base64.RawStdEncoding.EncodeToString(b)
	})
}

// Hex returns n random bytes as lowercase hex.
func Hex(n int) (string, error) {
	return fromEncoder(n, func(b []byte) string {
		return hex.EncodeToString(b)
	})
}

// Bytes returns n cryptographically secure random bytes.
func Bytes(n int) ([]byte, error) {
	if n <= 0 {
		return nil, fmt.Errorf("length must be > 0")
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %w", err)
	}
	return b, nil
}

// String returns n securely-generated characters using a custom charset.
// Charset length must be >= 2.
func String(n int, charset string) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("length must be > 0")
	}
	lenCharset := len(charset)
	if lenCharset < 2 {
		return "", fmt.Errorf("charset must have at least 2 characters")
	}

	out := make([]byte, n)
	// Rejection sampling threshold to avoid modulo bias
	// Accept only bytes < max, where max = 256 - (256 % lenCharset)
	max := 256 - (256 % lenCharset)
	buf := make([]byte, 1)

	for i := range n {
		for {
			if _, err := rand.Read(buf); err != nil {
				return "", fmt.Errorf("rand.Read failed: %w", err)
			}
			v := int(buf[0])
			if v < max {
				out[i] = charset[v%lenCharset]
				break
			}
			// otherwise discard and retry
		}
	}

	return string(out), nil
}

// ---- internal helper ----

// fromEncoder abstracts the pattern: random bytes → encode → string.
func fromEncoder(n int, enc func([]byte) string) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", err
	}
	return enc(b), nil
}
