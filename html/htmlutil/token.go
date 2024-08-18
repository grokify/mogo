package htmlutil

import (
	"fmt"
	"strings"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/type/stringsutil"
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

// Table returns a `[][]string` representing table data as text.
// Currently assumes input tokens represent one table and there are no nested tables.
// The output can be used with `github.com/grokify/gocharts/data/table`.
func (tokens Tokens) Table() [][]string {
	return newTableFromTokens(tokens)
}

func (tokens Tokens) Tokenizer() *html.Tokenizer {
	return NewTokenizerBytes([]byte(tokens.String()))
}

func ParseLink(tokens ...html.Token) (href string, desc string, err error) {
	if len(tokens) < 3 {
		return "", "", fmt.Errorf("less than 3 tokens, token count [%d]", len(tokens))
	}
	href, err = TokenAttribute(tokens[0], AttributeHref)
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
		"data":     t.Data,
		"string":   t.String()}
}

// MatchLeft matches the supplied token with the tokens in the set.
// only the attributes in the set need to match for a `true` result.
// One one set token need to match for success.
func (tokens Tokens) MatchLeft(tok html.Token, attrValMatchinfo *stringsutil.MatchInfo) bool {
	for _, tokFilterTry := range tokens {
		if TokenMatchLeft(tokFilterTry, tok, attrValMatchinfo) {
			return true
		}
	}
	return false
}

// TokenMatchLeft returns true if the token matches the token filter.
func TokenMatchLeft(tokFilter, tok html.Token, attrValMatchinfo *stringsutil.MatchInfo) bool {
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
		// since MatchInfo is being used as config against each attribute. If it is nil
		// set extact match with filter value.
		if attrValMatchinfo == nil {
			attrValMatchinfo = &stringsutil.MatchInfo{
				MatchType: stringsutil.MatchExact,
				String:    filAttr.Val,
			}
		}
		// since MatchInfo is being used as config against each attribute, populate
		// `MatchInfo` with `Attribute.Val`.
		if attrValMatchinfo.Regexp == nil && attrValMatchinfo.String == "" {
			attrValMatchinfo.String = filAttr.Val
		}
		// if !tokAttrs.Exists(filAttr, attrValMatchinfo) {
		// 	return false
		// }
		if tokAttrs.Index(filAttr, attrValMatchinfo) == -1 {
			return false
		}
	}
	return true
}

func (tokens Tokens) Subset(opts NextTokensOpts) Tokens {
	// func TokensSubset(startFilter, endFilter *TokenFilter, inclusive, greedy bool, toks []html.Token) []html.Token {
	// func TokensSubset(toks, start, end []html.Token, inclusive, greedy bool) []html.Token {
	subset := []html.Token{}
	if len(opts.StartFilter) == 0 && len(opts.EndFilter) == 0 {
		return tokens
	}
	matching := false
	if len(opts.StartFilter) == 0 {
		matching = true
	}

	for _, tok := range tokens {
		if opts.EndFilter.MatchLeft(tok, nil) {
			if matching && opts.InclusiveMatch {
				subset = append(subset, tok)
			}
			break
		} else if opts.StartFilter.MatchLeft(tok, opts.StartAttributeValueMatch) {
			if matching || opts.InclusiveMatch {
				subset = append(subset, tok)
			}
			matching = true
		} else if matching {
			subset = append(subset, tok)
		}
	}
	return subset
}
