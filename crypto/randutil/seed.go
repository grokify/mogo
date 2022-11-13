package randutil

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

// NewSeedInt64Crypto creates an `int64` seed value for `math/rand`.
// This is preferred over `NewSeedInt64Time()`.
// See: https://stackoverflow.com/a/54491783/1908967
func NewSeedInt64Crypto() (int64, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0, err
	}
	return int64(binary.LittleEndian.Uint64(b[:])), nil
}

// NewSeedInt64Time creates an `int64` seed value for `math/rand` based on
// `time.Now()`. This can have reduced entropy if used constantly throughout
// with shourl time differentials.
// See: https://stackoverflow.com/a/12321192/1908967
func NewSeedInt64Time() int64 {
	return time.Now().UnixNano()
}
