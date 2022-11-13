package mathutil

import (
	"math"
)

// some constants copied from https://github.com/golang/go/blob/master/src/math/bits.go
const (
	shift = 64 - 11 - 1
	bias  = 1023
	mask  = 0x7FF
)

// Round returns the nearest integer, rounding half away from zero.
// This function is available natively in Go 1.10
//
// Special cases are:
//
//	Round(±0) = ±0
//	Round(±Inf) = ±Inf
//	Round(NaN) = NaN
//
// This function is from the following gist by gdm85:
// https://gist.github.com/gdm85/44f648cc97bb3bf847f21c87e9d19b2d
func Round(x float64) float64 {
	// Round is a faster implementation of:
	//
	// func Round(x float64) float64 {
	//   t := Trunc(x)
	//   if Abs(x-t) >= 0.5 {
	//     return t + Copysign(1, x)
	//   }
	//   return t
	// }
	const (
		signMask = 1 << 63
		fracMask = 1<<shift - 1
		half     = 1 << (shift - 1)
		one      = bias << shift
	)

	bits := math.Float64bits(x)
	e := uint(bits>>shift) & mask
	if e < bias {
		// Round abs(x) < 1 including denormals.
		bits &= signMask // +-0
		if e == bias-1 {
			bits |= one // +-1
		}
	} else if e < bias+shift {
		// Round any abs(x) >= 1 containing a fractional component [0,1).
		//
		// Numbers with larger exponents are returned unchanged since they
		// must be either an integer, infinity, or NaN.
		e -= bias
		bits += half >> e
		bits &^= fracMask >> e
	}
	return math.Float64frombits(bits)
}

// Round is a rounding function for use before Go 1.11. The code
// was provided by David Vaini here:
// https://gist.github.com/DavidVaini/10308388
func RoundMore(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
