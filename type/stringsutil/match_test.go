package stringsutil

import (
	"strings"
	"testing"
)

var checkSuffixTests = []struct {
	full       string
	wantSuffix string
	gotFull    string
	gotPrefix  string
	gotSuffix  string
}{
	{"HelloWorld", "World", "HelloWorld", "Hello", "World"},
	{"HelloWorld", "Hello", "HelloWorld", "", ""},
}

func TestCheckSuffix(t *testing.T) {
	for _, tt := range checkSuffixTests {
		tryFull, tryPrefix, trySuffix := SuffixParse(tt.full, tt.wantSuffix)
		if tryFull != tt.gotFull || tryPrefix != tt.gotPrefix || trySuffix != tt.gotSuffix {
			t.Errorf("stringsutil.CheckSuffix() Error: want [%s,%s,%s], got [%s,%s,%s]",
				tt.gotFull, tt.gotPrefix, tt.gotSuffix, tryFull, tryPrefix, trySuffix)
		}
	}
}

var suffixMapTests = []struct {
	inputs           []string
	suffixes         []string
	prefixesString   string
	nonmatchesString string
}{
	{[]string{"HelloWorld", "HelloGalaxy", "ABC"}, []string{"World", "Galaxy", "Hello"}, "Hello", "ABC"},
}

func TestSuffixMap(t *testing.T) {
	for _, tt := range suffixMapTests {
		tryPrefixes, _, tryNonmatches := SuffixMap(tt.inputs, tt.suffixes)
		tryPrefixesString := strings.Join(tryPrefixes, ",")
		tryNonmatchesString := strings.Join(tryNonmatches, ",")
		if tryPrefixesString != tt.prefixesString || tryNonmatchesString != tt.nonmatchesString {
			t.Errorf("stringsutil.SuffixMap() Error: want [%s][%s], got [%s][%s]",
				tt.prefixesString, tt.nonmatchesString, tryPrefixesString, tryNonmatchesString)
		}
	}
}
