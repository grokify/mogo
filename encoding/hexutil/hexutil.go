package hexutil

import "encoding/hex"

// EncodeToString returns hex string using a supplied delimiter. The return
// value is lower case.
func EncodeToString(src []byte, delimit string) string {
	if delimit == "" {
		return hex.EncodeToString(src)
	}
	s := ""
	for i, b := range src {
		s += hex.EncodeToString([]byte{b})
		if i+1 < len(src) {
			s += delimit
		}
	}
	return s
}

// EncodeToStrings returns a slice of strings with leading 0s. The return
// values are lower case.
func EncodeToStrings(src []byte) []string {
	var s []string
	for _, b := range src {
		s = append(s, hex.EncodeToString([]byte{b}))
	}
	return s
}
