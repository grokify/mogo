package rsautil

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"hash"
)

type EncryptorOAEP struct {
	Md5           hash.Hash
	RsaPrivateKey *rsa.PrivateKey
	RsaPublicKey  *rsa.PublicKey
}

func NewEncryptorOAEP() EncryptorOAEP {
	enc := EncryptorOAEP{}
	enc.Md5 = md5.New()
	return enc
}

func (enc *EncryptorOAEP) DecryptOAEP(ciphertextBytes []byte, label []byte) ([]byte, error) {
	if enc.RsaPrivateKey == nil {
		err := errors.New("401: RSA Private Key Not Set")
		return []byte{}, err
	}
	return rsa.DecryptOAEP(enc.Md5, rand.Reader, enc.RsaPrivateKey, ciphertextBytes, label)
}

func (enc *EncryptorOAEP) DecryptOAEPBase64String(ciphertextBase64 string, label []byte) ([]byte, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return []byte(""), err
	}
	return enc.DecryptOAEP(ciphertextBytes, label)
}

func (enc *EncryptorOAEP) EncryptOAEP(plaintextBytes []byte, label []byte) ([]byte, error) {
	if enc.RsaPublicKey == nil {
		err := errors.New("402: RSA Public Key Not Set")
		return []byte{}, err
	}
	return rsa.EncryptOAEP(enc.Md5, rand.Reader, enc.RsaPublicKey, plaintextBytes, label)
}

func (enc *EncryptorOAEP) EncryptOAEPToBase64String(plaintextBytes []byte, label []byte) (string, error) {
	ciphertextBytes, err := enc.EncryptOAEP(plaintextBytes, label)
	if err != nil {
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(ciphertextBytes)
	return b64, nil
}
