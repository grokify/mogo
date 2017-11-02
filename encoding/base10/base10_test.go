// base10 supports Base10 encoding.
package base10

import (
	"testing"
)

var base10EncodeTests = []struct {
	v    []byte
	want int64
}{
	{[]byte("Hello World!"), int64(6810602152122453388)}}

func TestBase10Encode(t *testing.T) {
	for _, tt := range base10EncodeTests {
		enc := Encode(tt.v)

		if enc.Int64() != tt.want {
			t.Errorf("base10.Encode(%v): want %v, got %v", tt.v, tt.want, enc.Int64())
		}
	}
}
