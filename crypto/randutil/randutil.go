package randutil

import (
	crand "crypto/rand"
	"encoding/binary"
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

func Float64() (float64, error) {
	var b [8]byte
	if _, err := crand.Read(b[:]); err != nil {
		return 0, err
	}

	// Use the top 53 bits for uniform float64 precision (mimics math/rand.Float64()).
	// 1 << 53 is 9007199254740992, the number of representable values between 0 and 1 in float64.
	u := binary.BigEndian.Uint64(b[:]) >> 11 // 64 - 53 = 11
	return float64(u) / (1 << 53), nil
}

func MustFloat64() float64 {
	if f, err := Float64(); err != nil {
		panic(err)
	} else {
		return f
	}
}

// Intn returns a random number backed by `crypto/rand`.
func Intn(n uint) int {
	return mrand.New(NewCryptoRandSource()).Intn(int(n)) // #nosec G404 - `NewCryptoRandSource()` uses `crypto/rand`.
}

func Int64n(n uint) int64 {
	return int64(Intn(n))
}
