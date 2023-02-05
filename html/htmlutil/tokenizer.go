package htmlutil

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/grokify/mogo/type/stringsutil"
	"golang.org/x/net/html"
)

var (
	ErrTokenNotFound           = errors.New("token(s) not found")
	ErrTokenizerNotInitialized = errors.New("tokenizer not initialized")
)

func NewTokenizerBytes(b []byte) *html.Tokenizer {
	return html.NewTokenizer(bytes.NewReader(b))
}

func NewTokenizerFile(name string) (*html.Tokenizer, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return NewTokenizerBytes(b), nil
}

/*
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
*/

/*
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
*/

/*

func TokensBetweenNew(z *html.Tokenizer, skipErrors, inclusive bool, begin, end []html.Token) ([]html.Token, error) {
	begFilters := Tokens(begin)
	endFilter := Tokens(end)

	tokens := []html.Token{}
	tmsBegin, err := NextToken(z, skipErrors, begin...)
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
*/

func NextToken(z *html.Tokenizer, skipErrors bool, tokFilters ...html.Token) (html.Token, error) {
	opts := NextTokensOpts{
		SkipErrors:     skipErrors,
		IncludeChain:   false,
		InclusiveMatch: true,
		StartFilter:    []html.Token{},
		EndFilter:      tokFilters,
	}
	toks, err := NextTokens(z, opts)
	if err != nil {
		return html.Token{}, err
	} else if len(toks) == 0 {
		return html.Token{}, ErrTokenNotFound
	} else if len(toks) > 1 {
		panic("too many tokens (>1) found")
	}
	return toks[0], nil
}

type NextTokensOpts struct {
	SkipErrors               bool
	IncludeChain             bool
	InclusiveMatch           bool
	StartFilter              Tokens
	StartAttributeValueMatch *stringsutil.MatchInfo
	EndFilter                Tokens
}

func NextTokens(z *html.Tokenizer, opts NextTokensOpts) (Tokens, error) {
	// func NextTokens(z *html.Tokenizer, skipErrors, includeChain, includeMatch bool, start, end []html.Token) ([]html.Token, error) {
	matches := []html.Token{}
	if z == nil {
		return matches, ErrTokenizerNotInitialized
	}
	if opts.StartAttributeValueMatch == nil {
		opts.StartAttributeValueMatch = &stringsutil.MatchInfo{
			MatchType: stringsutil.MatchExact,
		}
	}
	foundStartToken := false
	if len(opts.StartFilter) == 0 {
		foundStartToken = true
	}

	for {
		ttThis := z.Next()
		switch ttThis {
		case html.ErrorToken:
			err := z.Err()
			if err == io.EOF {
				return matches, nil
			} else if !opts.SkipErrors {
				return matches, err
			}
		default:
			tok := z.Token()
			if foundStartToken {
				if opts.EndFilter.MatchLeft(tok, nil) {
					if opts.InclusiveMatch {
						matches = append(matches, tok)
					}
					return matches, nil
				} else if opts.IncludeChain {
					matches = append(matches, tok)
				}
			} else {
				if opts.StartFilter.MatchLeft(tok, opts.StartAttributeValueMatch) {
					if opts.InclusiveMatch {
						matches = append(matches, tok)
					}
					foundStartToken = true
				}
			}
		}
	}
}

// NextTextToken uses `NextTokensOpts` specifically for `SkipErrors`, `StartFilter`, and `StartAttributeValueMatch`.
func NextTextToken(z *html.Tokenizer, opts NextTokensOpts) (html.Token, error) {
	if z == nil {
		return html.Token{}, ErrTokenizerNotInitialized
	}
	if opts.StartAttributeValueMatch == nil {
		opts.StartAttributeValueMatch = &stringsutil.MatchInfo{
			MatchType: stringsutil.MatchExact,
		}
	}
	foundStartToken := false
	for {
		tokType := z.Next()
		tok := z.Token()
		if tokType == html.ErrorToken {
			err := z.Err()
			if err == io.EOF {
				return tok, ErrTokenNotFound
			} else if !opts.SkipErrors {
				return tok, err
			}
		} else if tokType == html.TextToken &&
			len(opts.StartFilter) == 0 || foundStartToken {
			return tok, nil
		} else if opts.StartFilter.MatchLeft(tok, opts.StartAttributeValueMatch) {
			foundStartToken = true
		}
	}
}

/*
// NextTokenMatch returns a string of matches. `includeMatch` is only used when `includeChain` is included.
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
*/

/*
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
*/

/*
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
*/

/*
func NextTextToken(z *html.Tokenizer, skipErrors bool, start []html.Token) (html.Token, error) {
	_, err := NextToken(z, skipErrors, start...)
	if err != nil {
		return html.Token{}, err
	}
	return NextToken(z, skipErrors, html.Token{Type: html.TextToken})
}
*/
