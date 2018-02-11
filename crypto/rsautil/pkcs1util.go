package rsautil

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
)

// SignRS256 signs data with rsa-sha256
func SignRS256(r *rsa.PrivateKey, data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r, crypto.SHA256, d)
}

// SignRS384 signs data with rsa-sha384
func SignRS384(r *rsa.PrivateKey, data []byte) ([]byte, error) {
	h := sha512.New384()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r, crypto.SHA384, d)
}

// SignRS512 signs data with rsa-sha512
func SignRS512(r *rsa.PrivateKey, data []byte) ([]byte, error) {
	h := sha512.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r, crypto.SHA512, d)
}

// SignRS512 signs data with rsa-sha512
func SignRS512Base64(r *rsa.PrivateKey, data []byte) (string, error) {
	sig, err := SignRS512(r, data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}
