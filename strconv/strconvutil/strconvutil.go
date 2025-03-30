package strconvutil

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/mogo/unicode/unicodeutil"
	"golang.org/x/exp/constraints"
)

func AnyToString(v any) string {
	if v == nil {
		return ""
	} else if valStr, ok := v.(string); ok {
		return valStr
	} else if valBool, ok := v.(bool); ok {
		return Btoa(valBool)
	} else if valInt, ok := v.(int); ok {
		return strconv.Itoa(valInt)
	} else if valInt, ok := v.(int8); ok {
		return strconv.Itoa(int(valInt))
	} else if valInt, ok := v.(int16); ok {
		return strconv.Itoa(int(valInt))
	} else if valInt, ok := v.(int32); ok {
		return strconv.Itoa(int(valInt))
	} else if valInt, ok := v.(int64); ok {
		return strconv.Itoa(int(valInt))
	} else if valInt, ok := v.(uint8); ok {
		return strconv.Itoa(int(valInt))
	} else if valInt, ok := v.(uint16); ok {
		return strconv.Itoa(int(valInt))
	} else if valInt, ok := v.(uint32); ok {
		return strconv.Itoa(int(valInt))
	} else if valInt, ok := v.(uint64); ok {
		return Itoa(valInt)
	} else if valTime, ok := v.(time.Time); ok {
		return valTime.Format(time.RFC3339)
	} else {
		return fmt.Sprintf("%v", v)
	}
}

// Float64ToString is a function type to define functions.
type Float64ToString func(float64) string

// Int64ToString is a function type to define functions.
type Int64ToString func(int64) string

// Btoa returns "true" or "false" according to the value of b.
func Btoa(b bool) string {
	return strconv.FormatBool(b)
}

// Commify takes an int64 and adds comma for every thousand
func Commify(n int64) string {
	// Stack Overflow: http://stackoverflow.com/users/1705598/icza
	// URL: https://stackoverflow.com/questions/13020308/how-to-fmt-printf-an-integer-with-thousands-comma
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
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
	// old
	// Stack Overflow: http://stackoverflow.com/users/1705598/icza
	// URL: http://stackoverflow.com/questions/13020308/how-to-fmt-printf-an-integer-with-thousands-comma
	/*
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
	*/
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

func FormatBoolMore(b bool, trueVal, falseVal string) string {
	if b {
		return trueVal
	} else {
		return falseVal
	}
}

func FormatDecimal[N constraints.Float | constraints.Integer](v N, precision int) string {
	if precision == 0 {
		return strconv.Itoa(int(v))
	} else if precision < 0 {
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	}
	return fmt.Sprintf(`%.`+strconv.Itoa(precision)+`f`, float64(v))
}

// ChangeToXoXPct converts a 1.0 == 100% based `float64` to a XoX percentage `float64`.
func ChangeToXoXPct(f float64) float64 {
	if f < 1.0 {
		return -1 * 100.0 * (1.0 - f)
	}
	return 100.0 * (f - 1.0)
}

// ChangeToFunnelPct converts a 1.0 == 100% based `float64` to a Funnel percentage `float64`.
func ChangeToFunnelPct(f float64) float64 { return f * 100.0 }

// Int64Len returns the length of an Int64 number.
func Int64Len(val int64) int {
	return len(fmt.Sprintf("%d", val))
}

// Int64Abbreviation returns integer abbreviations. For example,
// "1.5K", "15K", "150K", "1.5M", "15M", "150M".
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

// Itoa is like `strconv.Itoa()` with the additional functionality of converting `uint64` and accepting
// integer types natively via `constraints.Integer`.
func Itoa[E constraints.Integer](e E) string {
	return fmt.Sprintf("%d", e)
}

// UnquoteMore wraps `strconv.Unquote()` with additional functionality of allowing more chracters
// within single quotes.`
func UnquoteMore(s string) (string, error) {
	if len(s) < 2 {
		return s, nil
	}
	if string(s[0]) == unicodeutil.Apostrophe && string(s[len(s)-1]) == unicodeutil.Apostrophe {
		return s[1 : len(s)-1], nil
	}
	return strconv.Unquote(s)
}

// UnquoteMoreOrNot wraps `UnquoteMore()`
func UnquoteMoreOrNot(s string) string {
	if u, err := UnquoteMore(s); err == nil {
		return u
	} else {
		return s
	}
}
