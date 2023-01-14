package httputilmore

import (
	"testing"
)

var parseHTTPMethodTests = []struct {
	v          string
	wantString string
}{
	{" get ", "GET"},
	{"  PaTcH   ", "PATCH"},
}

func TestParseHTTPMethod(t *testing.T) {
	for _, tt := range parseHTTPMethodTests {
		canonicalStringTry, err := ParseHTTPMethodString(tt.v)

		if err != nil {
			t.Errorf("httputilmore.ParseHTTPMethodString(\"%s\") Error: [%s]",
				tt.v, err.Error())
		}
		if canonicalStringTry != tt.wantString {
			t.Errorf("httputilmore.ParseHTTPMethodString(\"%s\") Fail: want [%s] got [%s]",
				tt.v, tt.wantString, canonicalStringTry)
		}
	}
}
