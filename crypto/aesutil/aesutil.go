// Package aesutil provides authenticated AES-GCM encryption utilities,
// including writing and reading encrypted files.
//
// Ciphertext is a random 96-bit nonce followed by the GCM-sealed data and its
// 16-byte authentication tag, so decryption detects tampering and wrong keys.
// Keys must be 16, 24, or 32 bytes (AES-128, AES-192, or AES-256), and a key
// must not encrypt more than 2^32 messages, to keep random nonce collisions
// negligible.
//
// Before v0.75.0 this package used unauthenticated AES-CFB over a base64
// encoding of the plaintext. Ciphertext produced by those versions cannot be
// decrypted by this version.
package aesutil

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"os"
	"path"

	"github.com/btcsuite/btcd/btcutil/base58"
)

func EncryptAESBase58JSON(plainitem any, key []byte) ([]byte, error) {
	plaintext, err := json.Marshal(plainitem)
	if err != nil {
		return plaintext, err
	}
	return EncryptAESBase58(plaintext, key)
}

func DecryptAESBase58JSON(ciphertext []byte, key []byte, item any) error {
	plaintext, err := DecryptAESBase58(ciphertext, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(plaintext, item)
}

func EncryptAESBase58(plaintext []byte, key []byte) ([]byte, error) {
	bytes, err := EncryptAES(plaintext, key)
	if err != nil {
		return bytes, err
	}
	return []byte(base58.Encode(bytes)), nil
}

func DecryptAESBase58(ciphertext []byte, key []byte) ([]byte, error) {
	bytes := base58.Decode(string(ciphertext))
	return DecryptAES(bytes, key)
}

// EncryptAES encrypts plaintext with AES-GCM under key, returning the nonce
// followed by the sealed ciphertext and authentication tag.
func EncryptAES(plaintext []byte, key []byte) ([]byte, error) {
	aead, err := newAEAD(key)
	if err != nil {
		return nil, err
	}
	return aead.Seal(nil, nil, plaintext, nil), nil //nolint:gosec // G407: NewGCMWithRandomNonce requires a nil nonce and generates a random one
}

func DecryptAESBase64String(ciphertextBase64 string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return []byte{}, err
	}
	return DecryptAES(ciphertext, key)
}

// DecryptAES decrypts ciphertext produced by EncryptAES. It returns an error
// if the ciphertext is truncated, was modified, or was encrypted under a
// different key. The ciphertext slice is not modified.
func DecryptAES(ciphertext []byte, key []byte) ([]byte, error) {
	aead, err := newAEAD(key)
	if err != nil {
		return nil, err
	}
	return aead.Open(nil, nil, ciphertext, nil)
}

// newAEAD returns AES-GCM for key with random nonces that Seal prepends to
// the ciphertext and Open reads back.
func newAEAD(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCMWithRandomNonce(block)
}

func ReadFileAES(filename string, key []byte) ([]byte, error) {
	baFileEnc, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return DecryptAES(baFileEnc, key)
}

func WriteFileAES(filename string, baFileUnc []byte, perm os.FileMode, key []byte) error {
	baFileEnc, err := EncryptAES(baFileUnc, key)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, baFileEnc, perm) //nolint:gosec // G703: caller-supplied output path, as with os.WriteFile
}

func EncryptFileAES(filenameUnc string, filenameEnc string, perm os.FileMode, key []byte) error {
	baFileUnc, err := os.ReadFile(filenameUnc)
	if err != nil {
		return err
	}
	if len(filenameEnc) == 0 {
		filenameEnc = filenameUnc
	}
	return WriteFileAES(filenameEnc, baFileUnc, perm, key)
}

func EncryptDirectoryFilesAES(dirUnc string, dirEnc string, perm os.FileMode, key []byte) error {
	aFilesSrc, err := os.ReadDir(dirUnc)
	if err != nil {
		return err
	}
	for _, f := range aFilesSrc {
		if f.IsDir() {
			continue
		}
		fileUnc := f.Name()
		fileEnc := fileUnc + ".enc"
		pathUnc := path.Join(dirUnc, fileUnc)
		pathEnc := path.Join(dirEnc, fileEnc)
		err := EncryptFileAES(pathUnc, pathEnc, perm, key)
		if err != nil {
			return err
		}
	}
	return nil
}
