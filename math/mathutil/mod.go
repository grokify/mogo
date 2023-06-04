package mathutil

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Float | constraints.Integer
}

func ModInt[N Number](x, y N) int {
	return int(math.Mod(float64(x), float64(y)))
}

func ModInt64[N Number](x, y N) int64 {
	return int64(math.Mod(float64(x), float64(y)))
}

func ModPyInt[I constraints.Integer](a, b I) I {
	// https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
	// https://www.reddit.com/r/golang/comments/bnvik4/modulo_in_golang/
	return (a%b + b) % b
}

/*
func PyMod[I constraints.Integer](d, m I) I {
	// https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
	// https://www.reddit.com/r/golang/comments/bnvik4/modulo_in_golang/
	res := d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}
	return res
}

// ModInt provides a Python-like Modulo function.
func ModInt(a, b int) int {
	// https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
	// https://www.reddit.com/r/golang/comments/bnvik4/modulo_in_golang/
	return (a%b + b) % b
}

// ModInt64 provides a Python-like Modulo function.
func ModInt64(a, b int64) int64 {
	// https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
	// https://www.reddit.com/r/golang/comments/bnvik4/modulo_in_golang/
	return (a%b + b) % b
}
*/

// DivideInt64 performs integer division, returning
// a quotient and remainder.
func DivideInt64(dividend, divisor int64) (quotient, remainder int64) {
	// from https://stackoverflow.com/questions/43945675/division-with-returning-quotient-and-remainder
	quotient = dividend / divisor // integer division, decimals are truncated
	remainder = dividend % divisor
	return
}
