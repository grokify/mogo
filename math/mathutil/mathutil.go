package mathutil

// DivideInt64 performs integer division, returning
// a quotient and remainder.
func DivideInt64(dividend, divisor int64) (quotient, remainder int64) {
	// from https://stackoverflow.com/questions/43945675/division-with-returning-quotient-and-remainder
	quotient = dividend / divisor // integer division, decimals are truncated
	remainder = dividend % divisor
	return
}

// MinMaxInt32 returns min/max value given a slice of
// input values.
func MinMaxInt32(vals ...int32) (int32, int32) {
	min := int32(0)
	max := int32(0)
	for i, val := range vals {
		if i == 0 {
			min = val
			max = val
		} else {
			if val < min {
				min = val
			}
			if val > max {
				max = val
			}
		}
	}
	return min, max
}

// MinMaxUint returns min/max value given a slice of
// input values.
func MinMaxUint(vals ...uint) (uint, uint) {
	min := uint(0)
	max := uint(0)
	for i, val := range vals {
		if i == 0 {
			min = val
			max = val
		} else {
			if val < min {
				min = val
			}
			if val > max {
				max = val
			}
		}
	}
	return min, max
}
