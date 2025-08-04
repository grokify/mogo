package sqlutil

import (
	"testing"

	"github.com/grokify/mogo/text/foobar"
)

func TestBuildSQLXInsertSQLNamedParams(t *testing.T) {
	var buildSQLXInsertTests = []struct {
		vTableName   string
		vColumnNames []string
		want         string
	}{
		{"foobar", foobar.Vars(),
			`INSERT INTO foobar (bar,baz,corge,foo,fred,garply,grault,plugh,quux,qux,thud,waldo,xyzzy) VALUES (:bar,:baz,:corge,:foo,:fred,:garply,:grault,:plugh,:quux,:qux,:thud,:waldo,:xyzzy)`},
	}

	for _, tt := range buildSQLXInsertTests {
		got, err := BuildSQLXInsertSQLNamedParams(tt.vTableName, tt.vColumnNames)
		if err != nil {
			t.Errorf("err sqlutil.BuildSQLXInsert: err (%s)", err.Error())
		} else if got != tt.want {
			t.Errorf("mismatch sqlutil.BuildSQLXInsert(): want (%s) got (%s)", tt.want, got)
		}
	}
}
