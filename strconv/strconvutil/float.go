package strconvutil

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
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

func Ftoa[F constraints.Float](f F, prec int) string {
	return strconv.FormatFloat(float64(f), 'f', prec, 64)
}

func FormatFloat64ToIntStringFunnel[F constraints.Float](f F) string {
	return FormatFloat64ToAnyStringFunnel(float64(f), `%0.0f%%`)
}

// FormatFloat64ToAnyStringFunnel is used for funnels.
func FormatFloat64ToAnyStringFunnel(f float64, pattern string) string {
	return fmt.Sprintf(pattern, ChangeToFunnelPct(f))
}

func FormatFloat64ToIntString(f float64) string {
	return FormatFloat64ToAnyString(f, `%0.0f%%`)
}

// FormatFloat64ToAnyString is used for XoX growth.
func FormatFloat64ToAnyString(f float64, pattern string) string {
	return fmt.Sprintf(pattern, ChangeToXoXPct(f))
}
