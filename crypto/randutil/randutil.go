package randutil

import (
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"reflect"

	"github.com/grokify/mogo/math/mathutil"
	"github.com/grokify/mogo/type/number"
)

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
	if result, err := crand.Int(crand.Reader, max); err != nil {
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
// It supports any integer type, including named types, and full-width ranges.
func CryptoRandIntInRange[T number.Integer](min, max T) (T, error) {
	if min > max {
		return 0, errors.New("min must be <= max")
	}
	// mask is the largest value of T's width as an unsigned number. Masking
	// max-min recovers the true distance even when the subtraction wraps, as
	// it does for wide ranges of narrow signed types.
	mask := uint64(math.MaxUint64) >> (64 - reflect.TypeFor[T]().Bits())
	diff := uint64(max-min) & mask

	var b [8]byte
	for {
		if _, err := crand.Read(b[:]); err != nil {
			return 0, err
		}
		n := binary.BigEndian.Uint64(b[:]) & mask
		if diff == mask { // full range: every value of T is valid
			return min + T(n), nil
		}
		// Rejection sampling: accept only below the largest multiple of the
		// span that fits, so every result is equally likely.
		span := diff + 1
		if n < mask-(mask%span) {
			return min + T(n%span), nil
		}
	}
}

// NewMathRandCryptoSource generates a `*math/rand.Rand` using a random seed.
func NewMathRandCryptoSource() *rand.Rand {
	var seed int64
	err := binary.Read(crand.Reader, binary.LittleEndian, &seed)
	if err != nil {
		panic(fmt.Sprintf("failed to seed math/rand: %v", err))
	}
	return rand.New(rand.NewSource(seed)) // #nosec G404
}
