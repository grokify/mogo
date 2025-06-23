package randutil

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/grokify/mogo/math/mathutil"
	"github.com/grokify/mogo/type/number"
)

func Float64() (float64, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
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

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64
func Int63() int64 {
	return int64(Intn(mathutil.MaxInt63 + 1))
}

// Intn returns a cryptographically secure random int in [0, n). It panics if n <= 0.
func Intn(n int) int {
	if n <= 0 {
		panic("randutil: Intn requires n > 0")
	}
	max := big.NewInt(int64(n))
	if result, err := rand.Int(rand.Reader, max); err != nil {
		panic(fmt.Sprintf("randutil: failed to generate random number (%s)", err.Error()))
	} else {
		return int(result.Int64())
	}
}

/*
// Intn returns a random number backed by `crypto/rand`.
func Intn(n int) int {
	return mrand.New(NewCryptoRandSource()).Intn(n) // #nosec G404 - `NewCryptoRandSource()` uses `crypto/rand`.
}
*/

// CryptoRandIntInRange returns a cryptographically secure random integer in [min, max] (inclusive).
func CryptoRandIntInRange[T number.Integer](min, max T) (T, error) {
	if min > max {
		return 0, errors.New("min must be <= max")
	}
	span := uint64(max-min) + 1
	if span == 0 {
		return 0, errors.New("span overflow")
	}

	var nBytes int
	var maxUint uint64
	switch any(min).(type) {
	case int8, uint8:
		nBytes = 1
		maxUint = math.MaxUint8
	case int16, uint16:
		nBytes = 2
		maxUint = math.MaxUint16
	case int32, uint32:
		nBytes = 4
		maxUint = math.MaxUint32
	case int64, uint64, int, uint:
		nBytes = 8
		maxUint = math.MaxUint64
	default:
		return 0, errors.New("unsupported type")
	}

	limit := maxUint - (maxUint % span)
	b := make([]byte, nBytes)
	for {
		if _, err := rand.Read(b); err != nil {
			return 0, err
		}
		var n uint64
		switch nBytes {
		case 1:
			n = uint64(b[0])
		case 2:
			n = uint64(binary.BigEndian.Uint16(b))
		case 4:
			n = uint64(binary.BigEndian.Uint32(b))
		case 8:
			n = binary.BigEndian.Uint64(b)
		}
		if n < limit {
			return min + T(n%span), nil
		}
	}
}
