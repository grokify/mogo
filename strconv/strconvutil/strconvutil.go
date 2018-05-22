package strconvutil

import (
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
