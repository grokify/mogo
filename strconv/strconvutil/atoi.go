package strconvutil

import (
	"regexp"
	"strconv"
	"strings"
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

func CanonicalIntStringOrIgnore(s, comma, decimal string) string {
	try, err := AtoiMore(s, comma, decimal)
	if err != nil {
		return s
	}
	return strconv.Itoa(try)
}

func AtoiMore(s, comma, decimal string) (int, error) {
	if len(comma) > 0 {
		s = strings.Replace(s, comma, "", -1)
	}
	if len(decimal) > 0 && strings.Index(s, decimal) > -1 {
		s = regexp.MustCompile(regexp.QuoteMeta(decimal)+`.*$`).ReplaceAllString(s, "")
	}
	return strconv.Atoi(s)
}
