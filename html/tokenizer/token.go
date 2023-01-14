package tokenizer

import (
	"fmt"
	"strings"

	"github.com/grokify/mogo/errors/errorsutil"
	"golang.org/x/net/html"
)

type Tokens []html.Token

func (tokens Tokens) Maps() []map[string]string {
	maps := []map[string]string{}
	for _, tok := range tokens {
		maps = append(maps, TokenMap(tok))
	}
	return maps
}

func (tokens Tokens) String() string {
	toks := []string{}
	for _, tok := range tokens {
		toks = append(toks, tok.String())
	}
	return strings.Join(toks, "")
}

func ParseLink(tokens ...html.Token) (href string, desc string, err error) {
	if len(tokens) < 3 {
		return "", "", fmt.Errorf("less than 3 tokens, token count [%d]", len(tokens))
	}
	href, err = TokenAttribute(tokens[0], AttrHref)
	if err != nil {
		return href, "", errorsutil.Wrap(err,
			fmt.Sprintf("href not found in token [%s]",
				tokens[0].DataAtom))
	}
	desc = Tokens(tokens[1 : len(tokens)-1]).String()
	return
}

func TokenMap(t html.Token) map[string]string {
	return map[string]string{
		"type":     t.Type.String(),
		"dataAtom": t.DataAtom.String(),
		"string":   t.String()}
}

// Match matches the supplied token with the tokens in the set.
// only the attributes in the set need to match for a `true` result.
// One one set token need to match for success.
func (tokens Tokens) MatchLeft(tok html.Token) bool {
	for _, tokFilterTry := range tokens {
		if TokenMatchLeft(tokFilterTry, tok) {
			return true
		}
	}
	return false
}

// TokenMatchLeft returns true if the token matches the token filter.
func TokenMatchLeft(tokFilter, tok html.Token) bool {
	if tokFilter.Type != tok.Type {
		return false
	} else if tokFilter.DataAtom != tok.DataAtom {
		return false
	}
	if len(tokFilter.Attr) == 0 {
		return true
	}
	tokAttrs := Attributes(tok.Attr)
	for _, filAttr := range tokFilter.Attr {
		if !tokAttrs.Exists(filAttr) {
			return false
		}
	}
	return true
}
