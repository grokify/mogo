package config

import (
	"os"
	"testing"
)

var dmyhm2ParseTests = []struct {
	v1      string
	v2      string
	v3      string
	v4      string
	delimit string
	want    string
}{
	{"RINGCENTRAL_TOKEN", "ABC", "RINGCENTRAL_TOKEN_2", "def", "", "ABCdef"},
}

// TestJoinEnvNumbered ensures config.JoinEnvNumbered joins environment variables correctly.
func TestJoinEnvNumbered(t *testing.T) {
	for _, tt := range dmyhm2ParseTests {
		os.Setenv(tt.v1, tt.v2)
		os.Setenv(tt.v3, tt.v4)
		got := JoinEnvNumbered(tt.v1, tt.delimit, 2, true)
		if got != tt.want {
			t.Errorf("config.JoinEnvNumbered(\"%v\", \"\", 2, true) Mismatch: want %v, got %v", tt.v1, tt.want, got)
		}
	}
}
