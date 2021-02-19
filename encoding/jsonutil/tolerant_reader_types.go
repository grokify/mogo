package jsonutil

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/grokify/simplego/type/stringsutil"
)

// Bool implements a tolerant reader for `bool` type.
type Bool bool

var rxString = regexp.MustCompile(`^"(.*)"$`)

func (this *Bool) Value() bool {
	return bool(*this)
}

func (this *Bool) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	m := rxString.FindStringSubmatch(s)
	if len(m) == 2 {
		s = m[1]
	}
	*this = Bool(stringsutil.ToBool(s))
	return nil
}

// Int64 implements a tolerant reader for `int64` type.
type Int64 int64

func (this *Int64) Value() int64 {
	return int64(*this)
}

func (this *Int64) UnmarshalJSON(b []byte) error {
	*this = Int64(stringToInt64(string(b)))
	return nil
}

func stringToInt64(s string) int64 {
	s = strings.TrimSpace(s)
	m := rxString.FindStringSubmatch(s)
	if len(m) == 2 {
		s = strings.TrimSpace(m[1])
	}
	if len(s) == 0 {
		return 0
	}
	switch s {
	case "true":
		return 1
	case "false":
		return 0
	case "null":
		return 0
	}
	intVal, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return int64(intVal)
}
