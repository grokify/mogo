package ioutil

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/grokify/mogo/reflect/reflectutil"
)

var isReaderTests = []struct {
	v    any
	want bool
}{
	{bytes.NewReader([]byte{}), true},
	{strings.NewReader(" "), true},
}

func TestIsReader(t *testing.T) {
	for _, tt := range isReaderTests {
		isReader := IsReader(tt.v)
		if isReader != tt.want {
			t.Errorf("ioutil.IsReader(...) mismatch on type (%s) want (%s) got (%s)",
				reflectutil.NameOf(tt.v, true),
				strconv.FormatBool(tt.want), strconv.FormatBool(isReader))
		}
	}
}
