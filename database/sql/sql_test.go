package sql

import (
	"reflect"
	"strings"
	"testing"
)

var sliceToSQLsTests = []struct {
	sqlFormat    string
	valuesRaw    string
	valuesParsed []string
	want0        string
}{
	{"SELECT * FROM items WHERE name IN (%s)",
		"foo,bar,bax,qux", []string{"foo", "bar", "bax", "qux"},
		"SELECT * FROM items WHERE name IN ('foo','bar','bax','qux')"},
}

func TestSliceToSQLs(t *testing.T) {
	for _, tt := range sliceToSQLsTests {
		values := strings.Split(tt.valuesRaw, ",")
		if !reflect.DeepEqual(tt.valuesParsed, values) {
			t.Errorf("TestSliceToSQLs() panic, bad test: with [%v]", tt.valuesRaw)
		}
		sqls := BuildSQLsInStrings(tt.sqlFormat, values, MaxSQLLengthSOQL)
		if len(sqls) == 0 {
			t.Errorf("BuildSQLsInStrings() panic, bad test: with [%v] no results", tt.valuesRaw)
		}

		if tt.want0 != sqls[0] {
			t.Errorf("sqlBuildSQLsInStrings() Error: with [%v], want [%v], got [%v]",
				tt.valuesRaw, tt.want0, sqls[0])
		}
	}
}
