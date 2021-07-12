package tokenizer

import (
	"strings"

	"golang.org/x/net/html"
)

type Tokens []html.Token

func (tokens Tokens) String() string {
	toks := []string{}
	for _, tok := range tokens {
		toks = append(toks, tok.String())
	}
	return strings.Join(toks, "")
}
