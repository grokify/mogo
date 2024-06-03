package hexutil

import (
	"encoding/hex"
	"strings"
)

// EncodeToString returns hex string using a supplied seperator. The return
// value is lower case.
func EncodeToString(src []byte, sep string) string {
	if sep == "" {
		return hex.EncodeToString(src)
	}
	s := ""
	for i, b := range src {
		s += hex.EncodeToString([]byte{b})
		if i+1 < len(src) {
			s += sep
		}
	}
	return s
}

// EncodeToStrings returns a slice of strings with leading 0s. The `toUpper`
// param optionally returns upper case values.
func EncodeToStrings(src []byte, toUpper bool) []string {
	var s []string
	for _, b := range src {
		if toUpper {
			s = append(s, strings.ToUpper(hex.EncodeToString([]byte{b})))
		} else {
			s = append(s, hex.EncodeToString([]byte{b}))
		}
	}
	return s
}
