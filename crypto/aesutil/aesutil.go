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
	"io/ioutil"
	"os"
	"path"

	base58 "github.com/itchyny/base58-go"
)

func EncryptAesBase58Json(plainitem interface{}, key []byte, encoding *base58.Encoding) ([]byte, error) {
	plaintext, err := json.Marshal(plainitem)
	if err != nil {
		return plaintext, err
	}
	return EncryptAesBase58(plaintext, key, encoding)
}

func DecryptAesBase58Json(ciphertext []byte, key []byte, encoding *base58.Encoding, item interface{}) error {
	plaintext, err := DecryptAesBase58(ciphertext, key, encoding)
	if err != nil {
		return err
	}
	return json.Unmarshal(plaintext, item)
}

func EncryptAesBase58(plaintext []byte, key []byte, encoding *base58.Encoding) ([]byte, error) {
	bytes, err := EncryptAes(plaintext, key)
	if err != nil {
		return bytes, err
	}
	return encoding.Encode(bytes)
}

func DecryptAesBase58(ciphertext []byte, key []byte, encoding *base58.Encoding) ([]byte, error) {
	bytes, err := encoding.Decode(ciphertext)
	if err != nil {
		return bytes, err
	}
	return DecryptAes(bytes, key)
}

// EncryptAes provides a ciphertext byte array given a plaintext
// bytearray and key.
func EncryptAes(plaintext []byte, key []byte) ([]byte, error) {
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

func DecryptAesBase64String(ciphertextBase64 string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return []byte{}, err
	}
	return DecryptAes(ciphertext, key)
}

func DecryptAes(text []byte, key []byte) ([]byte, error) {
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

func ReadFileAes(filename string, key []byte) ([]byte, error) {
	baFileEnc, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return DecryptAes(baFileEnc, key)
}

func WriteFileAes(filename string, baFileUnc []byte, perm os.FileMode, key []byte) error {
	baFileEnc, err := EncryptAes(baFileUnc, key)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, baFileEnc, perm)
}

func EncryptFileAes(filenameUnc string, filenameEnc string, perm os.FileMode, key []byte) error {
	baFileUnc, err := ioutil.ReadFile(filenameUnc)
	if err != nil {
		return err
	}
	if len(filenameEnc) == 0 {
		filenameEnc = filenameUnc
	}
	return WriteFileAes(filenameEnc, baFileUnc, perm, key)
}

func EncryptDirectoryFilesAes(dirUnc string, dirEnc string, perm os.FileMode, key []byte) error {
	aFilesSrc, err := ioutil.ReadDir(dirUnc)
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
		err := EncryptFileAes(pathUnc, pathEnc, perm, key)
		if err != nil {
			return err
		}
	}
	return nil
}
