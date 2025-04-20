package htmlutil

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// newTableFromTokens returns a `[][]string` representing table data as text.
// Currently assumes input tokens represent one table and there are no nested tables.
// The output can be used with `github.com/grokify/gocharts/data/table`.
func newTableFromTokens(toks []html.Token) [][]string {
	inRow := false
	inCell := false
	rows := [][][]html.Token{}
	row := [][]html.Token{}
	cell := []html.Token{}
	for _, tok := range toks {
		if tok.DataAtom == atom.Tr {
			switch tok.Type {
			case html.StartTagToken:
				inRow = true
			case html.EndTagToken:
				rows = append(rows, row)
				row = [][]html.Token{}
				inRow = false
			}
			continue
		} else if tok.DataAtom == atom.Td || tok.DataAtom == atom.Th {
			switch tok.Type {
			case html.StartTagToken:
				inCell = true
			case html.EndTagToken:
				row = append(row, cell)
				cell = []html.Token{}
				inCell = false
			}
			continue
		} else if inRow && inCell {
			cell = append(cell, tok)
		}
	}
	rowsText := [][]string{}
	for _, row := range rows {
		rowText := []string{}
		for _, cell := range row {
			cellText := HTMLToText(Tokens(cell).String())
			cellText = strings.Join(strings.Fields(cellText), " ")
			rowText = append(rowText, cellText)
		}
		rowsText = append(rowsText, rowText)
	}

	return rowsText
}
