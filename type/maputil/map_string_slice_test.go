package maputil

import (
	"strings"
	"testing"
)

var mssTests = []struct {
	data           map[string][]string
	sortTestKey    string
	sortTestValues []string
}{
	{
		data: map[string][]string{
			"foo": {"foo", "bar", "baz"}},
		sortTestKey:    "foo",
		sortTestValues: []string{"bar", "baz", "foo"}},
}

func TestMapStringSlice(t *testing.T) {
	for _, tt := range mssTests {
		mss := MapStringSlice{}
		for key, vals := range tt.data {
			for _, val := range vals {
				mss.Add(key, val)
			}
		}
		mss.Sort(true)
		valsWant := strings.Join(tt.sortTestValues, ",")
		valsTry := strings.Join(mss[tt.sortTestKey], ",")
		if valsTry != valsWant {
			t.Errorf("maputil.MapStringSlice.Sort() Fail: want [%s], got [%s]",
				valsWant, valsTry)
		}
	}
}
