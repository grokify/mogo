package sha1util

import (
	"crypto/sha1"
	"crypto/sha512"
	"encoding/base32"
	"fmt"
	"io"
	"os"
)

func HashFile(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func Sum512d224Base32Bytes(b []byte, padChar rune) string {
	b28 := sha512.Sum512_224(b)
	return base32.StdEncoding.WithPadding(padChar).EncodeToString(b28[:])
}

func Sum512d224Base32String(s string, padChar rune) string {
	return Sum512d224Base32Bytes([]byte(s), padChar)
}

func Sum512d224Base32(r io.Reader, padChar rune) (string, error) {
	h := sha512.New512_224()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(padChar).EncodeToString(h.Sum([]byte{})), nil
}

func Sum512d224Base32File(name string, padChar rune) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return Sum512d224Base32(f, padChar)
}
