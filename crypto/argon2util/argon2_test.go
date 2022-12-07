package argon2util

import (
	"encoding/hex"
	"testing"
)

var hashTests = []struct {
	input string
	salt  string
	want  string
}{
	{
		"My Secret Password",
		"6368616e676520746869732070617373776f726420746f206120736563726574",
		"1XESZ8F6OGV200JDNXED8BAPU79UG5JA3KTJW7VMJCICX482S8"}}

func TestHash(t *testing.T) {
	for _, tt := range hashTests {
		saltBytes, err := hex.DecodeString(tt.salt)
		if err != nil {
			t.Errorf("TestRoundTrip cannot decode key: %v", err)
		}

		got := HashSimpleBase36([]byte(tt.input), saltBytes)

		if got != tt.want {
			t.Errorf("Argon2 Error: want [%v], got [%v]", tt.want, got)
		}
	}
}
