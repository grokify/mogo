package x509util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/grokify/mogo/os/osutil"
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

// GetRsaPrivateKeyForPkcs1PrivateKeyPath returns a *rsa.PrivateKey
// for a given PKCS#1 private key file path without a password
func GetRsaPrivateKeyForPkcs1PrivateKeyPath(prvKeyPKCS1Path string) (*rsa.PrivateKey, error) {
	isFile, err := osutil.IsFile(prvKeyPKCS1Path, true)
	if err != nil {
		return nil, err
	} else if !isFile {
		return nil, fmt.Errorf("filepath is not a file or is empty [%v]", prvKeyPKCS1Path)
	}

	prvKeyPkcs1Bytes, err := os.ReadFile(prvKeyPKCS1Path)
	if err != nil {
		return nil, err
	}

	return GetRsaPrivateKeyForPkcs1PrivateKeyBytes(prvKeyPkcs1Bytes)
}

func GetRsaPrivateKeyForPkcs1PrivateKeyBytes(prvKeyPkcs1Bytes []byte) (*rsa.PrivateKey, error) {
	block, rest := pem.Decode(prvKeyPkcs1Bytes)
	if len(rest) > 0 {
		return nil, fmt.Errorf("extra data included in key len: %v", len(rest))
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// GetRsaPrivateKeyForPkcs1PrivateKeyPathWithPassword returns a *rsa.PrivateKey
// for a given PKCS#1 private key file path and password
func GetRsaPrivateKeyForPkcs1PrivateKeyPathWithPassword(prvKeyPKCS1Path string, password []byte) (*rsa.PrivateKey, error) {
	isFile, err := osutil.IsFile(prvKeyPKCS1Path, true)
	if err != nil {
		return nil, err
	} else if !isFile {
		return nil, fmt.Errorf("filepath is not a file or is empty [%v]", prvKeyPKCS1Path)
	}

	prvKeyPkcs1BytesEnc, err := os.ReadFile(prvKeyPKCS1Path)
	if err != nil {
		return nil, err
	}

	return GetRsaPrivateKeyForPkcs1PrivateKeyBytesWithPassword(prvKeyPkcs1BytesEnc, password)
}

func GetRsaPrivateKeyForPkcs1PrivateKeyBytesWithPassword(prvKeyPkcs1BytesEnc []byte, password []byte) (*rsa.PrivateKey, error) {
	block, rest := pem.Decode(prvKeyPkcs1BytesEnc)
	if len(rest) > 0 {
		return nil, fmt.Errorf("extra data included in key len: %v", len(rest))
	}
	prvKeyBytes, err := x509.DecryptPEMBlock(block, password) // nolint:staticcheck
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PrivateKey(prvKeyBytes)
}

// GetRsaPublicKeyForPkcs8PublicKeyPath returns a *rsa.PublicKey
// for a given PKCS#8 public key file path.
func GetRsaPublicKeyForPkcs8PublicKeyPath(pubKeyPkcs8Path string) (*rsa.PublicKey, error) {
	var pubKey *rsa.PublicKey

	isFile, err := osutil.IsFile(pubKeyPkcs8Path, true)
	if err != nil {
		return pubKey, err
	} else if !isFile {
		return nil, fmt.Errorf("filepath is not a file or is empty [%v]", pubKeyPkcs8Path)
	}

	pubKeyPkcs8Bytes, err := os.ReadFile(pubKeyPkcs8Path)
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
