package randutil

import (
	crand "crypto/rand"
	"encoding/binary"
	"math"
	"math/rand"
	"time"

	"github.com/grokify/mogo/type/number"
)

// NewSeedInt64Crypto creates an `int64` seed value for `math/rand`.
// This is preferred over `NewSeedInt64Time()`.
// See: https://stackoverflow.com/a/54491783/1908967
func NewSeedInt64Crypto() (int64, error) {
	var b [8]byte
	_, err := crand.Read(b[:])
	if err != nil {
		return 0, err
	}
	v := binary.LittleEndian.Uint64(b[:]) & (1<<63 - 1)
	if v > uint64(math.MaxInt64) {
		v = uint64(math.MaxInt64)
	}
	return number.U64toi64(v)
	// return int64(v), nil
}

// NewSeedInt64Time creates an `int64` seed value for `math/rand`.
// This is preferred over `NewSeedInt64Crypto()`.
// See: https://stackoverflow.com/a/54491783/1908967
func NewSeedInt64Time() (int64, error) {
	return time.Now().UnixNano(), nil
}

// CryptoRandSource is a `crypto/rand` backed source that satisfies
// the `math/rand.Source` and `math/rand.Source64` interface definitions.
// It can be used as `r := rand.New(NewCryptoRandSource())`
// See: https://stackoverflow.com/a/35208651/1908967
type CryptoRandSource struct{}

func NewCryptoRandSource() rand.Source {
	return &CryptoRandSource{}
}

func (s *CryptoRandSource) Seed(seed int64) {
	// No-op: crypto/rand is always random, seeding is not supported.
}

func (s *CryptoRandSource) Int63() int64 {
	var v uint64
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		// Fallback or panic as appropriate for your application
		panic("crypto/rand failed: " + err.Error())
	}
	return int64(v & (1<<63 - 1))
}

func (s *CryptoRandSource) Uint64() uint64 {
	var b [8]byte
	_, err := crand.Read(b[:])
	if err != nil {
		// fallback to math/rand if crypto/rand fails
		x := mathrand.Int63()
		if x < 0 {
			return 0
		}
		return uint64(x)
	}
	return binary.LittleEndian.Uint64(b[:])
}

// #nosec G404 -- fallback to math/rand is acceptable here for non-crypto use
var mathrand = rand.New(rand.NewSource(time.Now().UnixNano()))
