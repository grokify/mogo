package x509util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

/*
 * Information on converting OpenSSH keys to OpenSSL PKCS1 and PKCS8 keys
 *
 * To decrypt OpenSSH Private Key to OpenSSL PKCS1 Private Key Format
 * openssl rsa -in id_rsa -out id_rsa.private.pkcs1
 *
 * To convert OpenSSH Public Key to OpenSSL PKCS8 Public Key Format
 * ssh-keygen -e -m PKCS8 -f id_rsa.pub > is_rsa.public.pkcs8
 *
 * Additional information on certificate formats:
 * https://www.netmeister.org/blog/ssh2pkcs8.html
 */

func GetRsaPrivateKeyForPkcs1PrivateKeyPath(prvKeyPKCS1Path string) (*rsa.PrivateKey, error) {
	var prvKey *rsa.PrivateKey

	prvKeyPkcs1Bytes, err := ioutil.ReadFile(prvKeyPKCS1Path)
	if err != nil {
		return prvKey, err
	}

	block, _ := pem.Decode(prvKeyPkcs1Bytes)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func GetRsaPrivateKeyForPkcs1PrivateKeyPathWithPassword(prvKeyPKCS1Path string, password []byte) (*rsa.PrivateKey, error) {
	var prvKey *rsa.PrivateKey

	prvKeyPkcs1BytesEnc, err := ioutil.ReadFile(prvKeyPKCS1Path)
	if err != nil {
		return prvKey, err
	}

	block, _ := pem.Decode(prvKeyPkcs1BytesEnc)
	prvKeyBytes, err := x509.DecryptPEMBlock(block, password)
	if err != nil {
		return prvKey, err
	}
	return x509.ParsePKCS1PrivateKey(prvKeyBytes)
}

func GetRsaPublicKeyForPkcs8PublicKeyPath(pubKeyPkcs8Path string) (*rsa.PublicKey, error) {
	var pubKey *rsa.PublicKey

	pubKeyPkcs8Bytes, err := ioutil.ReadFile(pubKeyPkcs8Path)
	if err != nil {
		return pubKey, err
	}

	block, _ := pem.Decode(pubKeyPkcs8Bytes)
	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return pubKey, err
	}

	pubKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return pubKey, errors.New("500: Cannot convert pub interface{} to *rsa.PublicKey")
	}
	return pubKey, nil
}
