package htmlutil

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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

func NewTokenFilter(tokenType html.TokenType, atoms ...atom.Atom) *TokenFilter {
	return &TokenFilter{
		TokenType: tokenType,
		AtomSet:   NewAtomSet(atoms...)}
}

func (tf *TokenFilter) Match(t html.Token) bool {
	if tf.AtomSet.Exists(t.DataAtom) &&
		t.Type == tf.TokenType {
		return true
	}
	return false
}

func TokensSubset(startFilter, endFilter *TokenFilter, inclusive, greedy bool, toks []html.Token) []html.Token {
	subset := []html.Token{}
	if startFilter == nil && endFilter == nil {
		return toks
	}
	matching := false
	if startFilter == nil {
		matching = true
	}
	for _, tok := range toks {
		if endFilter.Match(tok) {
			if matching && inclusive {
				subset = append(subset, tok)
			}
			break
		} else if startFilter.Match(tok) {
			if matching || inclusive {
				subset = append(subset, tok)
			}
			matching = true
		} else if matching {
			subset = append(subset, tok)
		}
	}
	return subset
}
