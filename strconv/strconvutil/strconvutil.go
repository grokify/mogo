package strconvutil

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// AtoiWithDefault is like Atoi but takes a default value
// which it returns in the event of a parse error.
func AtoiWithDefault(s string, def int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

// Commify takes an int64 and adds comma for every thousand
// Stack Overflow: http://stackoverflow.com/users/1705598/icza
// URL: http://stackoverflow.com/questions/13020308/how-to-fmt-printf-an-integer-with-thousands-comma
func Commify(n int64) string {
	in := strconv.FormatInt(n, 10)
	out := make([]byte, len(in)+(len(in)-2+int(in[0]/'0'))/3)
	if in[0] == '-' {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

var RxPlus = regexp.MustCompile(`^\+`)

func MustParseE164ToInt(s string) int {
	s = strings.TrimSpace(s)
	s = RxPlus.ReplaceAllString(s, "")
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func MustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func MustParseBool(s string) bool {
	parsed, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return parsed
}

// SliceStringToInt converts a slice of string integers.
func SliceStringToInt(strings []string) ([]int, error) {
	ints := []int{}
	for _, s := range strings {
		thisInt, err := strconv.Atoi(s)
		if err != nil {
			return ints, err
		}
		ints = append(ints, thisInt)
	}
	return ints, nil
}

// SliceStringToIntSort converts and sorts a slice of string integers.
func SliceStringToIntSort(strings []string) ([]int, error) {
	ints, err := SliceStringToInt(strings)
	if err != nil {
		return ints, err
	}
	intSlice := sort.IntSlice(ints)
	intSlice.Sort()
	return intSlice, nil
}

func FormatFloat64ToIntStringFunnel(v float64) string {
	return FormatFloat64ToAnyStringFunnel(v, `%0.0f%%`)
}

// FormatFloat64ToAnyStringFunnel is used for funnels.
func FormatFloat64ToAnyStringFunnel(v float64, pattern string) string {
	return fmt.Sprintf(pattern, ChangeToFunnelPct(v))
}

func FormatFloat64ToIntString(v float64) string {
	return FormatFloat64ToAnyString(v, `%0.0f%%`)
}

// FormatFloat64ToAnyString is used for XoX growth.
func FormatFloat64ToAnyString(v float64, pattern string) string {
	return fmt.Sprintf(pattern, ChangeToXoXPct(v))
}

// ChangeToXoXPct converts a 1.0 == 100% based `float64` to a
// XoX percentage `float64`.
func ChangeToXoXPct(v float64) float64 {
	if v < 1.0 {
		return -1 * 100.0 * (1.0 - v)
	}
	return 100.0 * (v - 1.0)
}

// ChangeToFunnelPct converts a 1.0 == 100% based `float64` to a
// Funnel percentage `float64`.
func ChangeToFunnelPct(v float64) float64 { return v * 100.0 }

// Int64Len returns the length of an Int64 number.
func Int64Len(val int64) int {
	return len(fmt.Sprintf("%d", val))
}

// Int64Abbreviation returns integer abbreviations.
func Int64Abbreviation(val int64) string {
	if val <= 999 {
		return strconv.Itoa(int(val))
	}
	valStr := fmt.Sprintf("%d", val)
	valLen := len(valStr)
	switch valLen {
	case 4:
		float := float64(val) / math.Pow10(valLen-1)
		return fmt.Sprintf("%.1fK", float)
	case 5:
		return fmt.Sprintf("%sK", valStr[0:2])
	case 6:
		return fmt.Sprintf("%sK", valStr[0:3])
	case 7:
		float := float64(val) / math.Pow10(valLen-1)
		return fmt.Sprintf("%.1fM", float)
	case 8:
		return fmt.Sprintf("%sM", valStr[0:2])
	case 9:
		return fmt.Sprintf("%sM", valStr[0:3])
	case 10:
		float := float64(val) / math.Pow10(valLen-1)
		return fmt.Sprintf("%.1fB", float)
	case 11:
		return fmt.Sprintf("%sB", valStr[0:2])
	case 12:
		return fmt.Sprintf("%sB", valStr[0:3])
	case 13:
		float := float64(val) / math.Pow10(valLen-1)
		return fmt.Sprintf("%.1fT", float)
	case 14:
		return fmt.Sprintf("%sT", valStr[0:2])
	case 15:
		return fmt.Sprintf("%sT", valStr[0:3])
	}
	return fmt.Sprintf("%d", val)
}
