package rsautil

import (
	"crypto/md5" // #nosec G501
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"hash"
	"os"
)

type CryptorOAEP struct {
	Hash          hash.Hash
	RsaPrivateKey *rsa.PrivateKey
	RsaPublicKey  *rsa.PublicKey
}

func NewCryptorOAEP() CryptorOAEP {
	return CryptorOAEP{Hash: md5.New()} // #nosec G401
}

func (enc *CryptorOAEP) DecryptOAEP(ciphertextBytes []byte, label []byte) ([]byte, error) {
	if enc.RsaPrivateKey == nil {
		return []byte{}, errors.New("401: RSA Private Key Not Set")
	}
	return rsa.DecryptOAEP(enc.Hash, rand.Reader, enc.RsaPrivateKey, ciphertextBytes, label)
}

func (enc *CryptorOAEP) DecryptOAEPBase64String(ciphertextBase64 string, label []byte) ([]byte, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return []byte(""), err
	}
	return enc.DecryptOAEP(ciphertextBytes, label)
}

func (enc *CryptorOAEP) DecryptOAEPBase64StringFromPath(ciphertextBase64Path string, label []byte) ([]byte, error) {
	ciphertextBase64Bytes, err := os.ReadFile(ciphertextBase64Path)
	if err != nil {
		return []byte(""), err
	}
	return enc.DecryptOAEPBase64String(string(ciphertextBase64Bytes), label)
}

func (enc *CryptorOAEP) EncryptOAEP(plaintextBytes []byte, label []byte) ([]byte, error) {
	if enc.RsaPublicKey == nil {
		err := errors.New("402: RSA Public Key Not Set")
		return []byte{}, err
	}
	return rsa.EncryptOAEP(enc.Hash, rand.Reader, enc.RsaPublicKey, plaintextBytes, label)
}

func (enc *CryptorOAEP) EncryptOAEPToBase64String(plaintextBytes []byte, label []byte) (string, error) {
	ciphertextBytes, err := enc.EncryptOAEP(plaintextBytes, label)
	if err != nil {
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(ciphertextBytes)
	return b64, nil
}
