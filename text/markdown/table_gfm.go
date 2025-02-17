package markdown

import (
	"strings"

	"github.com/grokify/mogo/type/slicesutil"
)

const (
	GFMTableSep      = "|"
	GFMTableSepStart = GFMTableSep + " "
	GFMTableSepMid   = " " + GFMTableSep + " "
	GFMTableSepEnd   = " " + GFMTableSep
)

func TableRowsToMarkdown(rows [][]string, newline string, esc, withHeader bool) string {
	var out string
	for i, row := range rows {
		md := TableRowToMarkdown(row, esc)
		out += md
		if i < len(rows)-1 {
			out += newline
		}
		if i == 0 && withHeader {
			out += TableSeparator(len(row))
			if i < len(rows)-1 {
				out += newline
			}
		}
	}
	return out
}

func TableRowToMarkdown(cells []string, esc bool) string {
	if !esc {
		return GFMTableSepStart + strings.Join(cells, GFMTableSepMid) + GFMTableSepEnd
	}
	new := []string{}
	for _, c := range cells {
		new = append(new, strings.ReplaceAll(c, `|`, `\|`))
	}
	return GFMTableSepStart + strings.Join(new, GFMTableSepMid) + GFMTableSepEnd
}

func TableSeparator(cellCount int) string {
	if cellCount == 0 {
		return ""
	}
	return GFMTableSepStart + strings.Join(slicesutil.NewWithDefault(cellCount, "-"), GFMTableSepMid) + GFMTableSepEnd
}
