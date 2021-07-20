package tokenizer

import (
	"errors"
	"fmt"
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

func ParseLink(tokens ...html.Token) (string, string, error) {
	if len(tokens) < 3 {
		return "", "", fmt.Errorf("less than 3 tokens [%d]", len(tokens))
	}
	href, err := TokenAttribute(tokens[0], AttrHref, true)
	if err != nil {
		return href, "", errors.New("href not found")
	}
	return href, Tokens(tokens[1 : len(tokens)-1]).String(), nil
}
