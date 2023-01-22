package bigint

import (
	"errors"
	"fmt"
	"math/big"
)

var ErrSetStringFailure = errors.New("string conversion to `*big.Int` failed")

// NewIntString creates a new `*big.Int` from an uint64.
func NewIntString(val string) (*big.Int, error) {
	i, ok := new(big.Int).SetString(val, 10)
	if !ok {
		return nil, ErrSetStringFailure
	}
	return i, nil
}

// NewIntHex converts a hex string to a `*big.Int`.
func NewIntHex(hexNumStr string) (*big.Int, error) {
	i, ok := new(big.Int).SetString(hexNumStr, 16)
	if !ok {
		return nil, ErrSetStringFailure
	}
	return i, nil
}

// NewIntUint64 creates a new `*big.Int` from an uint64.
func NewIntUint64(val uint64) *big.Int {
	return new(big.Int).SetUint64(val)
}

// IntToHex converts a `*big.Int` to a hex string.
func IntToHex(x *big.Int) string {
	return fmt.Sprintf("%x", x)
}

// Copy returns a copy of a `*big.Int`
func Copy(x *big.Int) *big.Int {
	return new(big.Int).SetBytes(x.Bytes())
}

// Div devides a by b and returns a new `*big.Int`
func Div(a, b *big.Int) *big.Int {
	return new(big.Int).Div(a, b)
}

// Equal checks if a == b.
func Equal(x, y *big.Int) bool {
	// https://tip.golang.org/src/math/big/alias_test.go
	return x.Cmp(y) == 0
}

// Mod performs `a mod n`
func Mod(a, n *big.Int) *big.Int {
	return new(big.Int).Mod(a, n)
}

// ModInt64 returns an int64 mod
func ModInt64(a, n int64) int64 {
	return new(big.Int).Mod(big.NewInt(n), big.NewInt(n)).Int64()
	// xBig := big.NewInt(x)
	// yBig := big.NewInt(y)
	// xBig.Mod(xBig, yBig)
	// return xBig.Int64()
}

// Pow is the power function for big ints.
func Pow(x, y *big.Int) *big.Int {
	// if y.String() == "1" {
	if Equal(y, big.NewInt(1)) {
		return Copy(x)
	} else if y.Sign() < 1 {
		return big.NewInt(1)
	} else if x.Sign() == 0 {
		return big.NewInt(0)
	}
	res := Copy(x)
	cyc := Copy(y)
	one := big.NewInt(1)
	for cyc.Cmp(one) > 0 {
		res = res.Mul(res, x)
		cyc = cyc.Sub(cyc, one)
	}
	return res
}

func PowInt64(x, y int64) int64 {
	return Pow(big.NewInt(x), big.NewInt(y)).Int64()
}

func IntToBaseXString(baseX, val int) string {
	return big.NewInt(int64(val)).Text(baseX)
}
