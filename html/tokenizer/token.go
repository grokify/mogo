package tokenizer

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
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

func ParseLink(tokens ...html.Token) (href string, desc string, err error) {
	if len(tokens) < 3 {
		return "", "", fmt.Errorf("less than 3 tokens, token count [%d]", len(tokens))
	}
	href, err = TokenAttribute(tokens[0], AttrHref)
	if err != nil {
		return href, "", errors.Wrap(err,
			fmt.Sprintf("href not found in token [%s]",
				tokens[0].DataAtom))
	}
	desc = Tokens(tokens[1 : len(tokens)-1]).String()
	return
}
