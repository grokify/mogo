package mathutil

// DivideInt64 performs integer division, returning
// a quotient and remainder.
func DivideInt64(dividend, divisor int64) (quotient, remainder int64) {
	quotient = dividend / divisor // integer division, decimals are truncated
	remainder = dividend % divisor
	return
}

// from https://stackoverflow.com/questions/43945675/division-with-returning-quotient-and-remainder

// MinMaxInt32 returns min/max value given two
// input values.
func MinMaxInt32(a, b int32) (int32, int32) {
	if b < a {
		return b, a
	}
	return a, b
}
