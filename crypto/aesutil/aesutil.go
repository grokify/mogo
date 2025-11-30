// aesutil provides AES crypto utilities including writing and reading
// AES encrypted files
package aesutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
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

// EncryptAes provides a ciphertext byte array given a plaintext bytearray and key.
func EncryptAES(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintextBase64 := base64.StdEncoding.EncodeToString(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(plaintextBase64))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintextBase64))
	return ciphertext, nil
}

func DecryptAESBase64String(ciphertextBase64 string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return []byte{}, err
	}
	return DecryptAES(ciphertext, key)
}

func DecryptAES(text []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("400: ciphertext is too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	plaintext, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return plaintext, nil
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
	return os.WriteFile(filename, baFileEnc, perm)
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
