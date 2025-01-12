// aesecb provides ECB (Electronic Codebook) mode for AES. Do not use this approach
// for new implementations due to weak security. This is provided for compatibility
// with previous implementations only.
package aesecb

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func DecryptBase64(v, key string) (string, error) {
	if decoded, err := base64.StdEncoding.DecodeString(v); err != nil {
		return "", err
	} else if block, err := aes.NewCipher([]byte(key)); err != nil {
		return "", err
	} else if len(decoded)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	} else {
		decrypted := make([]byte, len(decoded))
		decryptECB(block, decrypted, decoded)

		decrypted, err = removePKCS7Padding(decrypted, aes.BlockSize)
		if err != nil {
			return "", err
		} else {
			return string(decrypted), nil
		}
	}
}

// decryptECB performs the ECB mode decryption
func decryptECB(block cipher.Block, dst, src []byte) {
	bs := block.BlockSize()
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
}

// removePKCS7Padding removes padding from the decrypted data.
func removePKCS7Padding(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("invalid padding size")
	}

	paddingLen := int(data[len(data)-1])
	if paddingLen == 0 || paddingLen > blockSize {
		return nil, errors.New("invalid padding")
	}

	for _, padByte := range data[len(data)-paddingLen:] {
		if int(padByte) != paddingLen {
			return nil, errors.New("invalid padding")
		}
	}

	return data[:len(data)-paddingLen], nil
}

// EncryptBase64 provides ECB (Electronic Codebook) mode for AES. Do not use this approach
// for new implementations due to weak security. This is provided for compatibility
// with previous implementations only.
func EncryptBase64(v, key string) (string, error) {
	paddedInput := addPKCS5Padding([]byte(v), aes.BlockSize)

	if block, err := aes.NewCipher([]byte(key)); err != nil {
		return "", err
	} else {
		encrypted := make([]byte, len(paddedInput))
		encryptECB(block, encrypted, paddedInput)
		return base64.StdEncoding.EncodeToString(encrypted), nil
	}
}

// encryptECB performs the ECB mode encryption.
func encryptECB(block cipher.Block, dst, src []byte) {
	bs := block.BlockSize()
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
}

func addPKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}
