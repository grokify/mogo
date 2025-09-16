package aesargon

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/argon2"
)

type Crypter struct {
	Password  string
	UseBase58 bool
}

func (c *Crypter) Encrypt(plaintext string) (string, error) {
	return Encrypt(plaintext, c.Password, c.UseBase58)
}

func (c *Crypter) Decrypt(ciphertext string) (string, error) {
	return Decrypt(ciphertext, c.Password, c.UseBase58)
}

// deriveKeyAndIVArgon2id derives AES key and IV from password using Argon2id
func deriveKeyAndIVArgon2id(password, salt []byte) (key, iv []byte) {
	const (
		time    = 3         // iterations
		memory  = 64 * 1024 // 64 MB
		threads = 4
		keyLen  = 32
		ivLen   = 16
	)

	hash := argon2.IDKey(password, salt, time, memory, threads, keyLen+ivLen)
	key = hash[:keyLen]
	iv = hash[keyLen:]
	return
}

// Encrypt plaintext using AES-256-CBC with Argon2id-derived key
// useBase58: true = Base58 output, false = Base64 output
func Encrypt(plaintext, password string, useBase58 bool) (string, error) {
	salt := make([]byte, 16) // 16-byte salt
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key, iv := deriveKeyAndIVArgon2id([]byte(password), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	padLen := aes.BlockSize - len([]byte(plaintext))%aes.BlockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	plainPadded := append([]byte(plaintext), padding...)

	ciphertext := make([]byte, len(plainPadded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plainPadded)

	final := append(salt, ciphertext...)

	if useBase58 {
		return base58.Encode(final), nil
	}
	return base64.StdEncoding.EncodeToString(final), nil
}

// Decrypt ciphertext using AES-256-CBC with Argon2id-derived key
// useBase58: true = decode Base58, false = decode Base64
func Decrypt(ciphertextStr, password string, useBase58 bool) (string, error) {
	var data []byte
	var err error
	if useBase58 {
		data = base58.Decode(ciphertextStr)
	} else {
		data, err = base64.StdEncoding.DecodeString(ciphertextStr)
		if err != nil {
			return "", err
		}
	}

	if len(data) < 16 {
		return "", errors.New("ciphertext too short")
	}

	salt := data[:16]
	ciphertext := data[16:]
	key, iv := deriveKeyAndIVArgon2id([]byte(password), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS7 padding
	padLen := int(plaintext[len(plaintext)-1])
	if padLen > aes.BlockSize || padLen == 0 {
		return "", errors.New("invalid padding")
	}
	plaintext = plaintext[:len(plaintext)-padLen]

	return string(plaintext), nil
}
