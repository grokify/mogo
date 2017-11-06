package secretboxutil

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

// SealBase32String seals a message and returns a base32 encoded string.
// Example secret key: secretKeyBytes, err := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
func SealBase32String(plaintext, secretKeyBytes []byte) (string, error) {
	ciphertext, err := SealBox(plaintext, secretKeyBytes)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(ciphertext), nil
}

// SealBox seals a message using a supplied secret key and random nonce
// which is appended to the message.
func SealBox(plaintext, secretKeyBytes []byte) ([]byte, error) {
	var secretKey [32]byte
	copy(secretKey[:], secretKeyBytes)

	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return []byte(""), err
	}

	return secretbox.Seal(nonce[:], plaintext, &nonce, &secretKey), nil
}

// OpenBase32String opens a base32 encoded message.
func OpenBase32String(ciphertext32 string, secretKeyBytes []byte) ([]byte, error) {
	ciphertext, err := base32.StdEncoding.DecodeString(ciphertext32)
	if err != nil {
		return ciphertext, err
	}
	return OpenBox(ciphertext, secretKeyBytes)
}

// OpenBox opens a message which is prefixed by a nonce.
func OpenBox(ciphertext []byte, secretKeyBytes []byte) ([]byte, error) {
	var secretKey [32]byte
	copy(secretKey[:], secretKeyBytes)

	var nonce [24]byte
	copy(nonce[:], ciphertext[:24])
	plaintext, ok := secretbox.Open(nil, ciphertext[24:], &nonce, &secretKey)
	if !ok {
		return []byte(""), fmt.Errorf("Cannot decrypt")
	}
	return plaintext, nil
}
