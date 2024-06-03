package hexutil

import "encoding/hex"

// EncodeToString returns hex string using a supplied delimiter. The return
// value is lower case.
func EncodeToString(b []byte, delimit string) string {
	if delimit == "" {
		return hex.EncodeToString(b)
	}
	s := ""
	for i, bi := range b {
		s += hex.EncodeToString([]byte{bi})
		if i+1 < len(b) {
			s += delimit
		}
	}
	return s
}

// EncodeToStrings returns a slice of strings with leading 0s. The return
// values are lower case.
func EncodeToStrings(b []byte) []string {
	var s []string
	for _, bi := range b {
		s = append(s, hex.EncodeToString([]byte{bi}))
	}
	return s
}
