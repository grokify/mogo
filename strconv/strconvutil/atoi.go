package strconvutil

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrValueIsNegative   = errors.New("value is negative")
	ErrValueIsOutOfRange = errors.New("value is out of range")
)

// AtoiOrDefault is like Atoi but takes a default value
// which it returns in the event of a parse error.
func AtoiOrDefault(s string, def int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

func Atoi32(s string) (int32, error) {
	i32, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(i32), nil
}

func Atoi16(s string) (int16, error) {
	i16, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(i16), nil
}

func Atoi8(s string) (int8, error) {
	i8, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(i8), nil
}

func Atou(s string) (uint, error) {
	if i, err := strconv.Atoi(s); err != nil {
		return 0, err
	} else if i < 0 {
		return 0, ErrValueIsNegative
	} else {
		return uint(i), nil
	}
}

func Atou16(s string) (uint16, error) {
	if i, err := strconv.Atoi(s); err != nil {
		return 0, err
	} else if i < 0 {
		return 0, ErrValueIsNegative
	} else if i > int(^uint16(0)) {
		return 0, ErrValueIsOutOfRange
	} else {
		return uint16(i), nil
	}
}

func Atou32(s string) (uint32, error) {
	if i, err := strconv.Atoi(s); err != nil {
		return 0, err
	} else if i < 0 {
		return 0, ErrValueIsNegative
	} else if i > int(^uint32(0)) {
		return 0, ErrValueIsOutOfRange
	} else {
		return uint32(i), nil
	}
}

func CanonicalIntStringOrIgnore(s, comma, decimal string) string {
	try, err := AtoiMore(s, comma, decimal)
	if err != nil {
		return s
	}
	return strconv.Itoa(try)
}

func AtoiFunc(funcStringToInt64 func(s string) (int, error), s string) (int, error) {
	if funcStringToInt64 != nil {
		return funcStringToInt64(s)
	} else {
		return strconv.Atoi(s)
	}
}

func AtoiMore(s, comma, decimal string) (int, error) {
	if len(comma) > 0 {
		s = strings.ReplaceAll(s, comma, "")
	}
	if len(decimal) > 0 && strings.Contains(s, decimal) {
		s = regexp.MustCompile(regexp.QuoteMeta(decimal)+`.*$`).ReplaceAllString(s, "")
	}
	return strconv.Atoi(s)
}

func AtoiMoreFunc(comma, decimal string) func(s string) (int, error) {
	return func(s string) (int, error) {
		return AtoiMore(s, comma, decimal)
	}
}

// ParseUint8 safely parses a string into a uint8 value.
// It returns an error if the string is not a valid number
// or if the number is out of the uint8 range.
func ParseUint8(s string, base int) (uint8, error) {
	if u64, err := strconv.ParseUint(s, base, 8); err != nil {
		return 0, err
	} else if u64 > 255 {
		return 0, fmt.Errorf("value out of range for uint8 (%s), base (%d)", s, base)
	} else {
		return uint8(u64), nil
	}
}
