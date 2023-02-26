package jsonutil

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/grokify/mogo/reflect/reflectutil"
	"github.com/grokify/mogo/type/stringsutil"
)

// Bool implements a tolerant reader for `bool` type.
type Bool bool

var rxString = regexp.MustCompile(`^"(.*)"$`)

func (b *Bool) Value() bool {
	return bool(*b)
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	s := strings.TrimSpace(string(data))
	m := rxString.FindStringSubmatch(s)
	if len(m) == 2 {
		s = m[1]
	}
	*b = Bool(stringsutil.ToBool(s))
	return nil
}

// Int64 implements a tolerant reader for `int64` type.
type Int64 int64

func (i64 *Int64) Value() int64 {
	return int64(*i64)
}

func (i64 *Int64) UnmarshalJSON(data []byte) error {
	*i64 = Int64(stringToInt64(string(data)))
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

// String implements a tolerant reader for `string` type, including flexible input conversion to strings from numbers and integers.
type String string

func (s *String) UnmarshalJSON(data []byte) error {
	var a any

	err := json.Unmarshal(data, &a)
	if err != nil {
		panic(err)
	}

	if str, ok := a.(string); ok {
		*s = String(str)
		return nil
	} else if fl, ok := a.(float64); ok {
		*s = String(fmt.Sprintf("%g", fl))
		return nil
	} else if bl, ok := a.(bool); ok {
		if bl {
			*s = String("true")
		} else {
			*s = String("false")
		}
		return nil
	}

	return fmt.Errorf("json: cannot unmarshal %s into Go type github.com/mogo/encoding/jsonutil.String", reflectutil.NameOf(a, false))
}
