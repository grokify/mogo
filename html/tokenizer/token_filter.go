package tokenizer

import (
	"golang.org/x/net/html"
)

type TokenFilters []TokenFilter

func (filters TokenFilters) ByTokenType(tt html.TokenType) []TokenFilter {
	fils := []TokenFilter{}
	for _, fil := range filters {
		if fil.TokenType == tt {
			fils = append(fils, fil)
		}
	}
	return fils
}

// find next <tr> or </table>
type TokenFilter struct {
	TokenType html.TokenType
	AtomSet   AtomSet
}
