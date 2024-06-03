package hexutil

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

var encodeToStringTests = []struct {
	v    []byte
	sep  string
	want string
}{
	{[]byte{0x00, 0x00}, "", "0000"},
	{[]byte{0x00, 0x00}, " ", "00 00"},
	{[]byte{0x00, 0x00}, ":", "00:00"},
	{[]byte{0x00, 0x01}, "", "0001"},
	{[]byte{0x00, 0x01}, " ", "00 01"},
	{[]byte{0x00, 0x01}, ":", "00:01"},
	{[]byte{0x00, 0x01, 0x0a}, " ", "00 01 0a"},
	{[]byte{0x00, 0x01, 0x0a}, ":", "00:01:0a"},
	{[]byte{0xde, 0xad, 0xbe, 0xef}, "", "deadbeef"},
	{[]byte{0xde, 0xad, 0xbe, 0xef}, " ", "de ad be ef"},
	{[]byte{0xde, 0xad, 0xbe, 0xef}, ":", "de:ad:be:ef"},
	{[]byte{0xfe, 0xed, 0xbe, 0xef}, "", "feedbeef"},
	{[]byte{0xfe, 0xed, 0xbe, 0xef}, " ", "fe ed be ef"},
	{[]byte{0xfe, 0xed, 0xbe, 0xef}, ":", "fe:ed:be:ef"},
}

func TestEncodeToString(t *testing.T) {
	for _, tt := range encodeToStringTests {
		try := EncodeToString(tt.v, tt.sep)
		if try != tt.want {
			t.Errorf("hexutil.EncodeToString(\"%v\",\"%s\") Mismattch: want [%s], got [%s]", tt.v, tt.sep, tt.want, try)
		}
	}
}

var encodeToStringsTests = []struct {
	v       []byte
	toUpper bool
	want    []string
}{
	{[]byte{0xde, 0xad, 0xbe, 0xef}, false, []string{"de", "ad", "be", "ef"}},
	{[]byte{0xde, 0xad, 0xbe, 0xef}, true, []string{"DE", "AD", "BE", "EF"}},
	{[]byte{0xfe, 0xed, 0xbe, 0xef}, false, []string{"fe", "ed", "be", "ef"}},
	{[]byte{0xfe, 0xed, 0xbe, 0xef}, true, []string{"FE", "ED", "BE", "EF"}},
}

func TestEncodeToStrings(t *testing.T) {
	for _, tt := range encodeToStringsTests {
		try := EncodeToStrings(tt.v, tt.toUpper)
		if !slices.Equal(try, tt.want) {
			t.Errorf("hexutil.EncodeToStrings(\"%v\", %s) Mismattch: want [%s], got [%s]",
				tt.v, strconv.FormatBool(tt.toUpper), strings.Join(tt.want, " "), strings.Join(try, " "))
		}
	}
}
