package mathutil

func DivMod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

// from https://stackoverflow.com/questions/43945675/division-with-returning-quotient-and-remainder

func MinMaxInt32(a, b int32) (int32, int32) {
	if b < a {
		return b, a
	}
	return a, b
}
