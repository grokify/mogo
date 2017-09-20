package strconvutil

import (
	"strconv"
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
