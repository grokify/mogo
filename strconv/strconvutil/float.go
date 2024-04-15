package strconvutil

import (
	"strconv"
	"strings"
)

func Atof(s string, bitSizeIs32 bool) (float64, error) {
	if bitSizeIs32 {
		return strconv.ParseFloat(s, 32)
	} else {
		return strconv.ParseFloat(s, 64)
	}
}

func AtofMore(s string, bitSizeIs32 bool, comma string) (float64, error) {
	if len(comma) > 0 {
		s = strings.Replace(s, comma, "", -1)
	}
	if bitSizeIs32 {
		return strconv.ParseFloat(s, 32)
	} else {
		return strconv.ParseFloat(s, 64)
	}
}

func AtofFunc(funcStringToFloat func(s string) (float64, error), s string) (float64, error) {
	if funcStringToFloat != nil {
		return funcStringToFloat(s)
	} else {
		return Atof(s, false)
	}
}
