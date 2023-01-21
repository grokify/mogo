package reflectutil

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/grokify/mogo/data/currencyutil"
)

var nameOfTests = []struct {
	v                   any
	typename            string
	typenameWithPkgPath string
}{
	{float64(33.2), "float64", "float64"},
	{int(1), "int", "int"},
	{"abc", "string", "string"},
	{map[string]string{}, "", ""},
	{http.Client{}, "Client", "net/http.Client"},
	{url.Values{}, "Values", "net/url.Values"},
	{&url.Values{}, "*Values", "*net/url.Values"},
	{currencyutil.Amount{}, "Amount", "github.com/grokify/mogo/data/currencyutil.Amount"},
	{&currencyutil.Amount{}, "*Amount", "*github.com/grokify/mogo/data/currencyutil.Amount"},
}

// TestNameOf ensures `reflectutil.NameOf()`` returns correct values.
func TestNameOf(t *testing.T) {
	for _, tt := range nameOfTests {
		typeNameTry := NameOf(tt.v, false)
		if typeNameTry != tt.typename {
			t.Errorf("reflectutil.NameOf(\"%v\", false) Mismatch: want [%s], got [%s]", tt.v, typeNameTry, tt.typename)
		}
		typeNamePkgPathTry := NameOf(tt.v, true)
		if typeNamePkgPathTry != tt.typenameWithPkgPath {
			t.Errorf("reflectutil.NameOf(\"%v\", true) Mismatch: want [%s], got [%s]", tt.v, typeNamePkgPathTry, tt.typenameWithPkgPath)
		}
	}
}
