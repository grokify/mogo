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

func Sum512d224Base32Bytes(b []byte, padding rune) string {
	b28 := sha512.Sum512_224(b)
	return base32.StdEncoding.WithPadding(padding).EncodeToString(b28[:])
}

func Sum512d224Base32String(s string, padding rune) string {
	return Sum512d224Base32Bytes([]byte(s), padding)
}

func Sum512d224Base32(r io.Reader, padding rune) (string, error) {
	h := sha512.New512_224()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(padding).EncodeToString(h.Sum([]byte{})), nil
}

func Sum512d224Base32File(name string, padding rune) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return Sum512d224Base32(f, padding)
}
