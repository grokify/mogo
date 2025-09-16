// aesopenssl implements AES that is compatible with `aes-256-cbc -salt`.
package aesopenssl

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5" // #nosec G501 // used for key expansion and doesn't require collision resistance
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/btcsuite/btcutil/base58"
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

// deriveKeyAndIV mimics OpenSSL's EVP_BytesToKey
func deriveKeyAndIV(password []byte, salt []byte, keyLen, ivLen int) (key, iv []byte) {
	var d, dI []byte
	for len(d) < keyLen+ivLen {
		h := md5.New() // #nosec G401 // used for key expansion and doesn't require collision resistance
		h.Write(dI)
		h.Write(password)
		if salt != nil {
			h.Write(salt)
		}
		dI = h.Sum(nil)
		d = append(d, dI...)
	}
	return d[:keyLen], d[keyLen : keyLen+ivLen]
}

/*
// Encrypt is compatible with `aes-256-cbc -salt`.
// echo -n "mysecretstring" | openssl enc -aes-256-cbc -a -salt -pass pass:MyPassword
// -aes-256-cbc → encryption algorithm
// -a → base64 output (easier to handle in text)
// -pass pass:MyPassword → password (replace with your own)
func Encrypt(plaintext, password string) (string, error) {
	salt := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	key, iv := deriveKeyAndIV([]byte(password), salt, 32, 16) // AES-256-CBC

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

	final := append([]byte("Salted__"), salt...)
	final = append(final, ciphertext...)

	return base64.StdEncoding.EncodeToString(final), nil
}

// Decrypt is compatible with `aes-256-cbc -salt`.
// echo "ENCRYPTED_STRING_HERE" | openssl enc -aes-256-cbc -a -d -pass pass:MyPassword
func Decrypt(ciphertextBase64, password string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < 16 || string(ciphertext[:8]) != "Salted__" {
		return "", errors.New("invalid ciphertext")
	}
	salt := ciphertext[8:16]
	key, iv := deriveKeyAndIV([]byte(password), salt, 32, 16)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext = ciphertext[16:]
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// remove PKCS7 padding
	padLen := int(plaintext[len(plaintext)-1])
	if padLen > aes.BlockSize || padLen == 0 {
		return "", errors.New("invalid padding")
	}
	plaintext = plaintext[:len(plaintext)-padLen]

	return string(plaintext), nil
}
*/

// Encrypt encrypts plaintext and returns Base58 or Base64 encoded ciphertext
func Encrypt(plaintext, password string, useBase58 bool) (string, error) {
	salt := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	key, iv := deriveKeyAndIV([]byte(password), salt, 32, 16) // AES-256-CBC

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

	final := append([]byte("Salted__"), salt...)
	final = append(final, ciphertext...)

	if useBase58 {
		return base58.Encode(final), nil
	}
	return base64.StdEncoding.EncodeToString(final), nil
}

// Decrypt decrypts Base58 or Base64 encoded ciphertext
func Decrypt(ciphertextStr, password string, useBase58 bool) (string, error) {
	var ciphertext []byte
	var err error
	if useBase58 {
		ciphertext = base58.Decode(ciphertextStr)
	} else {
		ciphertext, err = base64.StdEncoding.DecodeString(ciphertextStr)
		if err != nil {
			return "", err
		}
	}

	if len(ciphertext) < 16 || string(ciphertext[:8]) != "Salted__" {
		return "", errors.New("invalid ciphertext")
	}

	salt := ciphertext[8:16]
	key, iv := deriveKeyAndIV([]byte(password), salt, 32, 16)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext = ciphertext[16:]
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// remove PKCS7 padding
	padLen := int(plaintext[len(plaintext)-1])
	if padLen > aes.BlockSize || padLen == 0 {
		return "", errors.New("invalid padding")
	}
	plaintext = plaintext[:len(plaintext)-padLen]

	return string(plaintext), nil
}
