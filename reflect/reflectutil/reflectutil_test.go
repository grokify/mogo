package reflectutil

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/grokify/mogo/database/datasource"
	"github.com/grokify/mogo/time/timeutil"
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
	{timeutil.TimeMore{}, "TimeMore", "github.com/grokify/mogo/time/timeutil.TimeMore"},
	{&timeutil.TimeMore{}, "*TimeMore", "*github.com/grokify/mogo/time/timeutil.TimeMore"},
}

// TestNameOf ensures `reflectutil.NameOf()` returns correct values.
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

var fieldTagValueTests = []struct {
	v         any
	fieldName string
	tagName   string
	tagValue  string
}{
	{datasource.DataSource{}, "Driver", "json", "driver"},
}

func TestFieldTagValue(t *testing.T) {
	for _, tt := range fieldTagValueTests {
		tagValue, err := FieldTagValue(tt.v, tt.fieldName, tt.tagName)
		if err != nil {
			t.Errorf("reflectutil.FieldTagValue(...) error (%s) want (%s)", err.Error(), tt.tagValue)
		}
		if tagValue != tt.tagValue {
			t.Errorf("reflectutil.FieldTagValue(...) mismatch: want (%s) got (%s)", tt.tagValue, tagValue)
		}
	}
}
