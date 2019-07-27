package bigutil

import (
	"fmt"
	"math/big"
)

// NewIntString creates a new `*big.Int` from an uint64.
func NewIntString(val string) *big.Int {
	i := new(big.Int)
	i.SetString(val, 10)
	return i
}

// NewIntUint64 creates a new `*big.Int` from an uint64.
func NewIntUint64(val uint64) *big.Int {
	i := new(big.Int)
	i.SetUint64(val)
	return i
}

// HexToInt converts a hex string to a `*big.Int`.
func HexToInt(hexNumStr string) *big.Int {
	i := new(big.Int)
	i.SetString(hexNumStr, 16)
	return i
}

// IntToHex converts a `*big.Int` to a hex string.
func IntToHex(n *big.Int) string {
	return fmt.Sprintf("%x", n)
}

// ModInt performs `a mod n`
func ModInt(a, n *big.Int) *big.Int {
	amodn := new(big.Int)
	return amodn.Mod(a, n)
}

// IsEqualInt checks if a == b.
func IsEqualInt(a, b *big.Int) bool {
	return a.String() == b.String()
}

// CopyInt returns a copy of a `*big.Int`
func CopyInt(i *big.Int) *big.Int {
	newInt := new(big.Int)
	newInt.SetString(i.String(), 10)
	return newInt
}

// PowInt is the power function for big ints.
func PowInt(x *big.Int, y *big.Int) *big.Int {
	if y.Sign() < 1 {
		return big.NewInt(1)
	} else if x.Sign() == 0 {
		return big.NewInt(0)
	}
	res := CopyInt(x)
	cyc := CopyInt(y)
	one := big.NewInt(1)
	for {
		if cyc.Cmp(one) < 1 {
			break
		}
		res = res.Mul(res, x)
		cyc = cyc.Sub(cyc, big.NewInt(1))
	}
	return res
}
