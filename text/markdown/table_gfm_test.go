package markdown

import (
	"strings"
	"testing"
)

func TestTableRowsToMarkdown(t *testing.T) {
	var markdownFormatTests = []struct {
		columns []string
		rows    [][]string
	}{
		{[]string{"a", "a"}, [][]string{{"bb", "bbb"}, {"cccc", "ccccc"}}},
	}

	for _, tt := range markdownFormatTests {
		allRows := [][]string{tt.columns}
		allRows = append(allRows, tt.rows...)

		md := TableRowsToMarkdown(allRows, "\n", true, true)

		lines := strings.Split(md, "\n")

		length := 0
		for i, l := range lines {
			if i == 0 {
				length = len(l)
			} else if len(l) != length {
				t.Errorf("Table.Markdown() Mismatch: first line length (%d) line (%d) length (%d)",
					length, i, len(l))
			}
		}
	}
}
