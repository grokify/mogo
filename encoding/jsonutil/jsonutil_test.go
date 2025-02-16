package jsonutil

import (
	"testing"
)

var marshalSimpleTests = []struct {
	v      map[string]any
	prefix string
	indent string
	want   string
}{
	{map[string]any{"foo": "bar"}, "", "", "{\"foo\":\"bar\"}"},
	{map[string]any{"foo": 123}, "", "", "{\"foo\":123}"},
}

func TestMarshalSimple(t *testing.T) {
	for _, tt := range marshalSimpleTests {
		try, err := MarshalSimple(tt.v, tt.prefix, tt.indent)
		if err != nil {
			t.Errorf("jsonutil.MarshalSimple: err (%s)", err.Error())
		}
		if string(try) != tt.want {
			t.Errorf("jsonutil.MarshalSimple: want (%s) got (%s)", tt.want, string(try))
		}
	}
}

var unmarshalMSITests = []struct {
	v   map[string]any
	foo string
}{
	{map[string]any{"foo": "bar"}, "bar"},
}

type testUnmarshalMSIStruct struct {
	Foo string `json:"foo"`
}

func TestUnmarshalMSI(t *testing.T) {
	for _, tt := range unmarshalMSITests {
		try := &testUnmarshalMSIStruct{}
		if err := UnmarshalMSI(tt.v, try); err != nil {
			t.Errorf("jsonutil.UnmarshalMSI: err (%s)", err.Error())
		} else if try.Foo != tt.foo {
			t.Errorf("jsonutil.UnmarshalMSI: want (%s) got (%s)", tt.foo, try.Foo)
		}
	}
}

var validateQuickTests = []struct {
	v     []byte
	fuzzy bool
	want  bool
}{
	{[]byte(""), true, false},
	{[]byte(""), false, false},
	{[]byte(`{"foo":"bar"}`), true, true},
	{[]byte(`     {"foo":"bar"}    `), true, true},
	{[]byte(`{"foo":"bar"}`), false, true},
	{[]byte(`<"foo":"bar">`), false, false},
	{[]byte(`["foo","bar"]`), false, true},
}

func TestValidateQuick(t *testing.T) {
	for _, tt := range validateQuickTests {
		try := ValidateQuick(tt.v, tt.fuzzy)
		if try != tt.want {
			t.Errorf("jsonutil.ValidateQuick(%v,%v): want (%v) got (%v)", tt.v, tt.fuzzy, tt.want, try)
		}
	}
}
