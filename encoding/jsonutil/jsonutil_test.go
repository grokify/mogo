package jsonutil

import (
	"testing"
)

var unmarshalMSITests = []struct {
	v   map[string]interface{}
	foo string
}{
	{map[string]interface{}{"foo": "bar"}, "bar"},
}

type testUnmarshalMSIStruct struct {
	Foo string `json:"foo"`
}

func TestUnmarshalMSI(t *testing.T) {
	for _, tt := range unmarshalMSITests {
		try := &testUnmarshalMSIStruct{}
		err := UnmarshalMSI(tt.v, try)
		if err != nil {
			t.Errorf("jsonutil.UnmarshalMSI: err [%s]", err.Error())
		}
		if try.Foo != tt.foo {
			t.Errorf("jsonutil.UnmarshalMSI: want [%s] got [%s]", tt.foo, try.Foo)
		}
	}
}
