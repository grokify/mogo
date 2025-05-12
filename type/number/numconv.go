package number

import (
	"errors"
	"math"

	"github.com/grokify/mogo/errors/errorsutil"
	"golang.org/x/exp/constraints"
)

var ErrOverflow = errors.New("integer overflow")

func Itoi8(i int) (int8, error) {
	if i < math.MinInt8 || i > math.MaxInt8 {
		return 0, errorsutil.Wrapf(ErrOverflow, "int value (%d) overflows int8", i)
	} else {
		return int8(i), nil
	}
}

func Itoi16(i int) (int16, error) {
	if i < math.MinInt16 || i > math.MaxInt16 {
		return 0, errorsutil.Wrapf(ErrOverflow, "int value (%d) overflows int16", i)
	} else {
		return int16(i), nil
	}
}

func Itoi32(i int) (int32, error) {
	if i < math.MinInt32 || i > math.MaxInt32 {
		return 0, errorsutil.Wrapf(ErrOverflow, "int value (%d) overflows int32", i)
	} else {
		return int32(i), nil
	}
}

func Itou(i int) (uint, error) {
	if i < 0 || uint(i) > ^uint(0) {
		return 0, errorsutil.Wrapf(ErrOverflow, "int value (%d) overflows uint", i)
	} else {
		return uint(i), nil
	}
}

func Itou16(i int) (uint16, error) {
	if i < 0 || i > int(^uint16(0)) {
		return 0, errorsutil.Wrapf(ErrOverflow, "int value (%d) overflows uint16", i)
	} else {
		return uint16(i), nil
	}
}

func Itou32(i int) (uint32, error) {
	if i < 0 || i > int(^uint32(0)) {
		return 0, errorsutil.Wrapf(ErrOverflow, "int value (%d) overflows uint32", i)
	} else {
		return uint32(i), nil
	}
}

func Itoi32s(ints []int) ([]int32, error) {
	var out []int32
	for _, val := range ints {
		if vout, err := Itoi32(val); err != nil {
			return out, err
		} else {
			out = append(out, vout)
		}
	}
	return out, nil
}

func Itous(ints []int) ([]uint, error) {
	var out []uint
	for _, val := range ints {
		if vout, err := Itou(val); err != nil {
			return out, err
		} else {
			out = append(out, vout)
		}
	}
	return out, nil
}

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

func U64toi64(u uint64) (int64, error) {
	if u > math.MaxInt64 {
		return 0, errorsutil.Wrapf(ErrOverflow, "uint64 value (%d) overflows int64", u)
	} else {
		return int64(u), nil
	}
}
