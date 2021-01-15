package mathutil

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

// DivideInt64 performs integer division, returning
// a quotient and remainder.
func DivideInt64(dividend, divisor int64) (quotient, remainder int64) {
	// from https://stackoverflow.com/questions/43945675/division-with-returning-quotient-and-remainder
	quotient = dividend / divisor // integer division, decimals are truncated
	remainder = dividend % divisor
	return
}
