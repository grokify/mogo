package maputil

import (
	"testing"
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
