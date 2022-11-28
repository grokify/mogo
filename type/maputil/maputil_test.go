package maputil

import (
	"testing"

	"golang.org/x/exp/slices"
)

var mapSSEqualTests = []struct {
	v    map[string]string
	v2   map[string]string
	want bool
}{
	{map[string]string{"KEY": "num", "VAL": "123"}, map[string]string{"KEY": "num", "VAL": "123"}, true},
	{map[string]string{"first": "foo", "second": "bar"}, map[string]string{"first": "foo", "second": "baz"}, false},
}

func TestMapSSEqual(t *testing.T) {
	for _, tt := range mapSSEqualTests {
		eq := MapSSEqual(tt.v, tt.v2)

		if eq != tt.want {
			t.Errorf("maputil.MapMSSEqual() Error: with [%v][%v], want [%v], got [%v]",
				tt.v, tt.v2, tt.want, eq)
		}
	}
}

var stringKeysExistTests = []struct {
	v          map[string]any
	keys       []string
	requireAll bool
	keysExist  bool
}{
	{map[string]any{"Foo": "num", "Bar": "123", "Baz": true}, []string{"Foo", "Bar", "Baz"}, false, true},
	{map[string]any{"Foo": "num", "Bar": "123", "Baz": true}, []string{"Foo", "Bar", "Baz"}, true, true},
	{map[string]any{"Foo": "num", "Bar": "123", "Baz": true}, []string{"Foo", "Bar"}, false, true},
	{map[string]any{"Foo": "num", "Bar": "123", "Baz": true}, []string{"Foo", "Bar"}, true, true},
	{map[string]any{"Foo": "num", "Bar": "123", "Baz": true}, []string{"Foo", "Bar", "Qux"}, false, true},
	{map[string]any{"Foo": "num", "Bar": "123", "Baz": true}, []string{"Foo", "Bar", "Qux"}, true, false},
	{map[string]any{"KEY": "num", "VAL": "123"}, []string{"KEY", "VAL"}, true, true},
	{map[string]any{"KEY": "num", "VAL": "123"}, []string{"KEYVAL"}, true, false},
}

func TestStringKeysExist(t *testing.T) {
	for _, tt := range stringKeysExistTests {
		keysExistTry := StringKeysExist(tt.v, tt.keys, tt.requireAll)
		if tt.keysExist != keysExistTry {
			t.Errorf("maputil.StringKeysExist() params [%v] keys [%v] reqAll [%v] want [%v] got [%v]",
				tt.v, tt.keys, tt.requireAll, tt.keysExist, keysExistTry)
		}
	}
}

var stringValuesTests = []struct {
	v    map[string]string
	want []string
}{
	{map[string]string{"1": "foo", "2": "bar", "3": "baz"}, []string{"bar", "baz", "foo"}},
}

func TestStringValues(t *testing.T) {
	for _, tt := range stringValuesTests {
		got := StringValues(tt.v)
		if !slices.Equal(tt.want, got) {
			t.Errorf("maputil.StringValues() params [%v] want [%v] got [%v]",
				tt.v, tt.want, got)
		}
	}
}
