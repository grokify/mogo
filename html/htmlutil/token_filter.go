package htmlutil

/*
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
*/
