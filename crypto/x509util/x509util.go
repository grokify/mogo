package x509util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"github.com/grokify/gotilla/io/ioutilmore"
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
 * http://help.globalscape.com/help/eft6/Server_SSH_Key_Formats.htm
 * https://www.openssl.org/docs/apps/rsa.html
 * https://burnz.wordpress.com/2007/12/14/ssh-convert-openssh-to-ssh2-and-vise-versa/
 * http://www.sysmic.org/dotclear/index.php?post/2010/03/24/Convert-keys-betweens-GnuPG%2C-OpenSsh-and-OpenSSL
 * https://support.aerofs.com/hc/en-us/articles/202868994-How-Do-I-Convert-My-SSL-Certificate-File-To-PEM-Format-
 * https://shanetully.com/2012/04/simple-public-key-encryption-with-rsa-and-openssl/
 * https://www.socketloop.com/tutorials/golang-saving-private-and-public-key-to-files
 */

func GetRsaPrivateKeyForPkcs1PrivateKeyPath(prvKeyPKCS1Path string) (*rsa.PrivateKey, error) {
	var prvKey *rsa.PrivateKey

	isFileGtZero, err := ioutilmore.IsFileWithSizeGtZero(prvKeyPKCS1Path)
	if err != nil {
		return prvKey, err
	} else if isFileGtZero == false {
		return prvKey, errors.New("400: key file path is zero size.")
	}

	prvKeyPkcs1Bytes, err := ioutil.ReadFile(prvKeyPKCS1Path)
	if err != nil {
		return prvKey, err
	}

	block, _ := pem.Decode(prvKeyPkcs1Bytes)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func GetRsaPrivateKeyForPkcs1PrivateKeyPathWithPassword(prvKeyPKCS1Path string, password []byte) (*rsa.PrivateKey, error) {
	var prvKey *rsa.PrivateKey

	isFileGtZero, err := ioutilmore.IsFileWithSizeGtZero(prvKeyPKCS1Path)
	if err != nil {
		return prvKey, err
	} else if isFileGtZero == false {
		return prvKey, errors.New("400: key file path is zero size.")
	}

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

	isFileGtZero, err := ioutilmore.IsFileWithSizeGtZero(pubKeyPkcs8Path)
	if err != nil {
		return pubKey, err
	} else if isFileGtZero == false {
		return pubKey, errors.New("400: key file path is zero size.")
	}

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
