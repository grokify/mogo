package reflectutil

import (
	"net/http"
	"net/url"
	"testing"
)

var sliceMinMaxTests = []struct {
	v        any
	typename string
}{
	{float64(33.2), "float64"},
	{int(1), "int"},
	{"abc", "string"},
	{http.Client{}, "Client"},
	{url.Values{}, "Values"},
}

// TestTypeName ensures reflectutil.TypeName returns correct values.
func TestTypeName(t *testing.T) {
	for _, tt := range sliceMinMaxTests {
		typeNameTry := TypeName(tt.v)
		if typeNameTry != tt.typename {
			t.Errorf("reflectutil.TypeName(\"%v\") Mismatch: want [%s], got [%s]", tt.v, typeNameTry, tt.typename)
		}
	}
}
