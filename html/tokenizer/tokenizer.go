package tokenizer

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var TokenNotFound = errors.New("token not found")

func NewTokenizerFile(filename string) (*html.Tokenizer, error) {
	htmlBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return html.NewTokenizer(bytes.NewReader(htmlBytes)), nil
}

// TokensBetweenAtom returns the tokens that represent the `innerHtml`
// between a start and end tag token.
func TokensBetweenAtom(z *html.Tokenizer, skipErrors, inclusive bool, htmlAtom atom.Atom) ([]html.Token, error) {
	return TokensBetween(z, skipErrors, inclusive,
		TokenFilters{{
			TokenType: html.StartTagToken,
			AtomSet:   NewAtomSet(htmlAtom)}},
		TokenFilters{{
			TokenType: html.EndTagToken,
			AtomSet:   NewAtomSet(htmlAtom)}})
}

func TokensBetween(z *html.Tokenizer, skipErrors, inclusive bool, begin, end TokenFilters) ([]html.Token, error) {
	tokens := []html.Token{}
	tmsBegin, err := NextTokenMatch(z, skipErrors, false, false, begin...)
	if err != nil {
		return tokens, err
	}
	if inclusive {
		tokens = append(tokens, tmsBegin...)
	}
	tokensChain, err := NextTokenMatch(z, skipErrors, true, inclusive, end...)
	if err != nil {
		return tokens, err
	}
	tokens = append(tokens, tokensChain...)
	return tokens, nil
}

// NextTokenMatch returns a string of matches. `includeMatch` is only used
// when `includeChain` is included.
func NextTokenMatch(z *html.Tokenizer, skipErrors, includeChain, includeMatch bool, filters ...TokenFilter) ([]html.Token, error) {
	matches := []html.Token{}
	if len(filters) == 0 {
		return matches, errors.New("no filters provided")
	}
	filtersMore := TokenFilters(filters)
	for {
		tt := z.Next()
		token := z.Token()
		if token.Type == html.ErrorToken {
			break
		}
		filtersForType := filtersMore.ByTokenType(tt)
		if len(filtersForType) > 0 {
			for _, filter := range filtersForType {
				if filter.AtomSet.Len() == 0 {
					if !includeChain || includeMatch {
						matches = append(matches, token)
					}
					return matches, nil
				} else if filter.AtomSet.Exists(token.DataAtom) {
					if !includeChain || includeMatch {
						matches = append(matches, token)
					}
					return matches, nil
				}
			}
		}
		if includeChain {
			matches = append(matches, token)
		}
	}
	return matches, nil
}

func NextStartToken(z *html.Tokenizer, skipErrors bool, htmlAtoms ...atom.Atom) (html.Token, error) {
	if len(htmlAtoms) == 0 {
		return html.Token{}, errors.New("no atoms requested")
	}
	atoms := NewAtomSet(htmlAtoms...)
	for {
		ttThis := z.Next()
		switch ttThis {
		case html.ErrorToken:
			if !skipErrors {
				return html.Token{}, z.Err()
			}
		case html.StartTagToken:
			tok := z.Token()
			if atoms.Exists(tok.DataAtom) {
				return tok, nil
			}
		}
	}
	return html.Token{},
		fmt.Errorf("token not found for [%s]", strings.Join(atoms.Names(), ","))
}

func NextTextToken(z *html.Tokenizer, skipErrors bool, htmlAtoms ...atom.Atom) (html.Token, error) {
	atoms := NewAtomSet(htmlAtoms...)
	for {
		tokType := z.Next()
		tok := z.Token()
		if !skipErrors && tokType == html.ErrorToken {
			return tok, z.Err()
		} else if atoms.Len() == 0 && tokType == html.TextToken {
			return tok, nil
		} else if atoms.Len() > 0 &&
			tokType == html.StartTagToken &&
			atoms.Exists(tok.DataAtom) {
			return NextTextToken(z, skipErrors)
		}
	}
	return html.Token{}, TokenNotFound
}
