// base10 supports Base10 encoding.
package csvutil

import (
	"testing"
)

var newReaderTests = []struct {
	filename   string
	separator  rune
	columns    []string
	colNameLen []int
}{
	{"testdata/simple.csv", ',', []string{"bazqux", "foobar"}, []int{6, 6}},
	{"testdata/utf8bom.csv", ',', []string{"foobar", "bazqux"}, []int{6, 6}},
}

func TestNewReader(t *testing.T) {
	for _, tt := range newReaderTests {
		csvReader, f, err := NewReader(tt.filename, tt.separator)
		if err != nil {
			t.Errorf("csvutil.NewReader(\"%s\",...): error [%s]",
				tt.filename, err.Error())
		}
		defer f.Close()
		line, err := csvReader.Read()
		if err != nil {
			t.Errorf("csvutil.NewReader(\"%s\",...): csvReader.Read() error [%s]",
				tt.filename, err.Error())
		}
		if len(line) == 0 {
			t.Errorf("csvutil.NewReader(\"%s\",...): line is empty",
				tt.filename)
		}
		colName1Try := line[0]
		if len(colName1Try) != len(tt.columns[0]) {
			t.Errorf("csvutil.NewReader(\"%s\",...): colName mismatch want [%d] got [%d]",
				tt.filename, len(tt.columns[0]), len(colName1Try))
		}
	}
}
