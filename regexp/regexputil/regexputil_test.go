package regexputil

import (
	"reflect"
	"regexp"
	"testing"
)

var findStringSubmatchNamedMapTests = []struct {
	v    string
	v2   string
	want map[string]string
}{
	{`(?P<KEY>\w+):(?P<VAL>\w+)`, "num:123", map[string]string{"KEY": "num", "VAL": "123"}},
	{`(?P<first>\w+).(?P<second>\w+)`, "foo.bar", map[string]string{"first": "foo", "second": "bar"}},
}

func TestFindStringSubmatchNamedMap(t *testing.T) {
	for _, tt := range findStringSubmatchNamedMapTests {
		rx := regexp.MustCompile(tt.v)
		resMss := FindStringSubmatchNamedMap(rx, tt.v2)

		eq := reflect.DeepEqual(tt.want, resMss)
		if !eq {
			t.Errorf("regepxutil.FindStringSubmatchNamedMap() Error: with [%v][%v], want [%v], got [%v]",
				tt.v, tt.v2, tt.want, resMss)
		}
	}
}
