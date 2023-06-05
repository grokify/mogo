package slicesutil

import "github.com/grokify/mogo/math/mathutil"

func ReverseIndex(n, i uint) uint {
	if i >= n {
		i = mathutil.ModPyInt(i, n)
	}
	return n - i - 1
}
