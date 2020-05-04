package mathutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func PrettyTicksPercent(estimatedTickCount int, low, high float64, sigDigits uint8) []float64 {
	if sigDigits == 0 {
		sigDigits = 1
	}
	lowInt := int64(low * 100 * math.Pow10(int(sigDigits)))
	highInt := int64(high * 100 * math.Pow10(int(sigDigits)))
	ticksInt := PrettyTicks(estimatedTickCount, lowInt, highInt)
	ticks := []float64{}
	for _, tickVal := range ticksInt {
		ticks = append(ticks, float64(tickVal)/100/math.Pow10(int(sigDigits)))
	}
	return ticks
}

// PrettyTicks returns a slice of integers that start
// lower and end higher than the supplied range. This
// is intended to be used for chart axis.
func PrettyTicks(estimatedTickCount int, low, high int64) []int64 {
	ticks := []int64{}
	if low > high {
		tmp := low
		low = high
		high = tmp
	}
	diffRaw := high - low
	tickSize := float64(diffRaw) / float64(estimatedTickCount)
	tickSizedRounded := FloorMostSignificant(int64(tickSize))
	lowFloor := FloorMostSignificant(int64(low))

	ticks = append(ticks, lowFloor)
	for ticks[len(ticks)-1] < high {
		ticks = append(ticks, ticks[len(ticks)-1]+tickSizedRounded)
	}
	return ticks
}

// FloorMostSignificant returns number with a single
// significant digit followed by zeros. The value returned
// is always lower than the supplied value.
func FloorMostSignificant(valOriginal int64) int64 {
	// see here for additional discussion
	// https://stackoverflow.com/questions/202302/
	if valOriginal == 0 {
		return 0
	}
	valPositive := valOriginal
	isNegative := false
	if valOriginal < 0 {
		valPositive = -1 * valOriginal
		isNegative = true
	}
	valStr := fmt.Sprintf("%d", valPositive)
	valLen := len(fmt.Sprintf("%d", valPositive))
	var final int64

	// Math power approach
	mostSig, err := strconv.Atoi(valStr[0:1])
	if err != nil {
		panic(errors.Wrap(err, "mathutil.FloorMostSignificant"))
	}
	if isNegative {
		final = -1 * int64(mostSig+1) * int64(math.Pow10(valLen-1))
	} else {
		final = int64(mostSig) * int64(math.Pow10(valLen-1))
	}
	// String approach
	if 1 == 0 {
		vals := make([]string, valLen)
		for i := 0; i < valLen; i++ {
			if i == 0 {
				vals[i] = valStr[0:1]
			} else {
				vals[i] = "0"
			}
		}
		intStr := strings.Join(vals, "")
		num, err := strconv.Atoi(intStr)
		if err != nil {
			panic(errors.Wrap(err, "mathutil.FloorMostSignificant"))
		}
		final = int64(num)
	}
	return final
}
