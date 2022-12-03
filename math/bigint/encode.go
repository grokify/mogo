package bigint

import (
	"errors"
	"fmt"
	"math/big"
)

var ErrBaseInvalid = errors.New("invalid base - must be between [2,62]")

// EncodeToString uses [0-9a-zA-Z] for base62.
func EncodeToString(base int, s []byte) (string, error) {
	if base < 2 || base > 62 {
		return "", ErrBaseInvalid
	}
	var i big.Int
	i.SetBytes(s)
	return i.Text(base), nil
}

// MustEncodeToString uses [0-9a-zA-Z] for base62.
func MustEncodeToString(base int, s []byte) string {
	var i big.Int
	i.SetBytes(s)
	return i.Text(base)
}

func DecodeString(base int, s string) ([]byte, error) {
	if base < 2 || base > 62 {
		return []byte{}, ErrBaseInvalid
	}
	var i big.Int
	_, ok := i.SetString(s, base)
	if !ok {
		return []byte{}, fmt.Errorf("cannot parse base(%d) (%q)", base, s)
	}
	return i.Bytes(), nil
}
