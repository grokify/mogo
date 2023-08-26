package slicesutil

import "github.com/grokify/mogo/math/mathutil"

// ReverseIndex returns the forward index value from the end of the string.
func ReverseIndex(n, i uint) uint {
	if i >= n {
		i = mathutil.ModPyInt(i, n)
	}
	return n - i - 1
}

// IndexValueOrDefault returns the value at the supplied index or a supplied default value.
func IndexValueOrDefault[E any](s []E, idx uint, def E) E {
	if int(idx) >= len(s) {
		return def
	}
	return s[idx]
}
