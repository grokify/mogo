package randutil

import (
	mrand "math/rand"
)

/*
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
		mrand:  mrand.New(s)} // #nosec G404
}

func (cr *CryptoRand) Intn(n int) (int, error) {
	i64, err := cr.Int64n(int64(n))
	return int(i64), err
}

func (cr *CryptoRand) MustIntn(n int) int {
	return int(cr.MustInt64n(int64(n)))
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
*/

// Intn returns a random number backed by `crypto/rand`.
func Intn(n uint) int {
	return mrand.New(NewCryptoRandSource()).Intn(int(n)) // #nosec G404 - `NewCryptoRandSource()` uses `crypto/rand`.
}

func Int64n(n uint) int64 {
	return int64(Intn(n))
}
