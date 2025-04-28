package number

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

/*
const maxExactInt = 1 << 53  // 9007199254740992
const minExactInt = -1 << 53 // -9007199254740992
*/

func Ntof64[N constraints.Integer | constraints.Float](val N) (float64, bool) {
	const maxExactInt = 1 << 53

	f := float64(val)

	switch any(val).(type) {
	case int64, int:
		if f >= -float64(maxExactInt) && f <= float64(maxExactInt) {
			return f, true
		}
	case uint64, uint, uintptr:
		if f <= float64(maxExactInt) {
			return f, true
		}
	case int8, int16, int32, uint8, uint16, uint32, float32, float64:
		return f, true
	}
	return f, false
}

func IntToUint32(x int) (uint32, error) {
	if x < 0 || x > int(^uint32(0)) {
		return 0, fmt.Errorf("int is out of range for uint32 (%d)", x)
	}
	return uint32(x), nil
}
