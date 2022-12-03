package hmacutil

import (
	"testing"
)

var hmacSHA256HexTests = []struct {
	key    string
	msg    string
	macHex string
	macB32 string
}{
	{"5eb63bbbe01eeed093cb22bb8f5acdc3", `{"foo":"bar"}`,
		"8e06998c7c691a2a11e91d2326ba18126e26140566cb66fcbd4e167c43c41b57",
		"RYDJTDD4NENCUEPJDURSNOQYCJXCMFAFM3FWN7F5JYLHYQ6EDNLQ===="},
}

func TestHMACSHA256Hex(t *testing.T) {
	for _, tt := range hmacSHA256HexTests {
		gotHex := HMACSHA256Hex([]byte(tt.key), []byte(tt.msg))
		if gotHex != tt.macHex {
			t.Errorf("hmacutil.HMACSHA256Hex(\"%s\", \"%s\"): want [%s], got [%s]",
				tt.key, tt.msg, tt.macHex, gotHex)
		}
		gotB32 := HMACSHA256Base32([]byte(tt.key), []byte(tt.msg))
		if gotB32 != tt.macB32 {
			t.Errorf("hmacutil.HMACSHA256Base32(\"%s\", \"%s\"): want [%s], got [%s]",
				tt.key, tt.msg, tt.macB32, gotB32)
		}
	}
}
