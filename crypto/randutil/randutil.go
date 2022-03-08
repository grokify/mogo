package randutil

import (
	"crypto/rand"
	"io"
	"math/big"
	mrand "math/rand"
	"time"
)

type CryptoRand struct {
	reader io.Reader
	mrand  *mrand.Rand
}

func NewCryptoRand(r io.Reader, s mrand.Source) CryptoRand {
	if r == nil {
		r = rand.Reader
	}
	if s == nil {
		s = mrand.NewSource(time.Now().Unix())
	}
	return CryptoRand{
		reader: r,
		mrand:  mrand.New(s)}
}

func (cr *CryptoRand) Intn(n int) (int, error) {
	bign, err := rand.Int(cr.reader, big.NewInt(int64(n)))
	if err != nil {
		return -1, err
	}
	return int(bign.Int64()), nil
}

func (cr *CryptoRand) MustIntn(n int) int {
	i, err := cr.Intn(n)
	if err != nil {
		return cr.mrand.Intn(n)
	}
	return i
}

func (cr *CryptoRand) Int64n(n int64) (int64, error) {
	bign, err := rand.Int(cr.reader, big.NewInt(n))
	if err != nil {
		return -1, err
	}
	return bign.Int64(), nil
}

func (cr *CryptoRand) MustInt64n(n int64) int64 {
	i64, err := cr.Int64n(n)
	if err != nil {
		return int64(cr.mrand.Intn(int(n)))
	}
	return i64
}
