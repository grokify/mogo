package tokenizer

import (
	"bytes"
	"errors"
	"io"
	"os"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var ErrTokenNotFound = errors.New("token(s) not found")

func NewTokenizerBytes(b []byte) *html.Tokenizer {
	return html.NewTokenizer(bytes.NewReader(b))
}

func NewTokenizerFile(name string) (*html.Tokenizer, error) {
	htmlBytes, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return NewTokenizerBytes(htmlBytes), nil
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

func NextToken(z *html.Tokenizer, skipErrors bool, tokFilters ...html.Token) (html.Token, error) {
	if z == nil {
		return html.Token{}, errors.New("tokenizer must be supplied")
	}
	tokFil := Tokens(tokFilters)
	for {
		ttThis := z.Next()
		switch ttThis {
		case html.ErrorToken:
			err := z.Err()
			if z.Err() == io.EOF {
				return html.Token{}, ErrTokenNotFound
			} else if !skipErrors {
				return html.Token{}, err
			}
		default:
			tok := z.Token()
			// if tok.DataAtom == atom.Img {
			//	fmtutil.PrintJSON(tok)
			// }
			if len(tokFilters) == 0 || tokFil.MatchLeft(tok) {
				return tok, nil
			}
		}
	}
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
			err := z.Err()
			if z.Err() == io.EOF {
				return html.Token{}, ErrTokenNotFound
			} else if !skipErrors {
				return html.Token{}, err
			}
		case html.StartTagToken:
			tok := z.Token()
			if atoms.Exists(tok.DataAtom) {
				return tok, nil
			}
		}
	}
}

func NextTextToken(z *html.Tokenizer, skipErrors bool, htmlAtoms ...atom.Atom) (html.Token, error) {
	atoms := NewAtomSet(htmlAtoms...)
	for {
		tokType := z.Next()
		tok := z.Token()
		if tokType == html.ErrorToken {
			err := z.Err()
			if err == io.EOF {
				return tok, ErrTokenNotFound
			} else if !skipErrors {
				return tok, err
			}
		} else if atoms.Len() == 0 && tokType == html.TextToken {
			return tok, nil
		} else if atoms.Len() > 0 &&
			tokType == html.StartTagToken &&
			atoms.Exists(tok.DataAtom) {
			return NextTextToken(z, skipErrors)
		}
	}
}
