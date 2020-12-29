package join

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grokify/simplego/type/stringsutil"
)

// JoinAny takes an array of interface{} and converts
// each value to a string using fmt.Sprintf("%v")
func JoinAny(a []interface{}, sep string) string {
	strs := []string{}
	for _, item := range a {
		strs = append(strs, fmt.Sprintf("%v", item))
	}
	return strings.Join(strs, sep)
}

func JoinInt(a []int, sep string) string {
	strs := []string{}
	for _, i := range a {
		strs = append(strs, strconv.Itoa(i))
	}
	return strings.Join(strs, sep)
}

func JoinCondenseTrimSpace(slice []string, sep string) string {
	return strings.Join(stringsutil.SliceTrimSpace(slice, true), sep)
}

func JoinQuoteMaxLength(slice []string, begQuote, endQuote, sep string, maxLength int) []string {
	words := []string{}
	curWords := []string{}
	curLength := 0
	for _, word := range slice {
		if maxLength > 0 && curLength+len(begQuote+word+endQuote+sep) > maxLength {
			words = append(words, strings.Join(curWords, sep))
			curWords = []string{}
			curLength = 0
		} else {
			curWords = append(curWords, begQuote+word+endQuote)
			curLength += len(begQuote + word + endQuote + sep)
		}
	}
	if len(curWords) > 0 {
		words = append(words, strings.Join(curWords, sep))
	}
	return words
}

type JoinMoreOpts struct {
	QuoteBegin string
	QuoteEnd   string
	Separator  string
	MaxLength  int
	TrimSpace  bool
	SkipEmpty  bool
}

func JoinMore(slice []string, cfg JoinMoreOpts) []string {
	lines := []string{}
	curWords := []string{}
	curLength := 0
	for _, word := range slice {
		if cfg.TrimSpace {
			word = strings.TrimSpace(word)
		}
		if cfg.SkipEmpty && len(word) == 0 {
			continue
		}
		if cfg.MaxLength > 0 {
			if curLength+len(cfg.QuoteBegin+word+cfg.QuoteEnd+cfg.Separator) > cfg.MaxLength {
				lines = append(lines, strings.Join(curWords, cfg.Separator))
				curWords = []string{}
				curLength = 0
			} else {
				curWords = append(curWords, cfg.QuoteBegin+word+cfg.QuoteEnd)
				curLength += len(cfg.QuoteBegin + word + cfg.QuoteEnd + cfg.Separator)
			}
		} else {
			curWords = append(curWords, cfg.QuoteBegin+word+cfg.QuoteEnd)
		}
	}
	if len(curWords) > 0 {
		lines = append(lines, strings.Join(curWords, cfg.Separator))
	}
	return lines
}

func JoinQuoteTrimSpaceSkipEmpty(slice []string, begQuote, endQuote, sep string) string {
	strs := JoinMore(
		slice,
		JoinMoreOpts{
			QuoteBegin: begQuote,
			QuoteEnd:   endQuote,
			Separator:  sep,
			MaxLength:  -1,
			TrimSpace:  true,
			SkipEmpty:  true})
	if len(strs) == 0 {
		return ""
	}
	return strs[0]
}

func JoinQuoteMaxLengthTrimSpaceSkipEmpty(slice []string, begQuote, endQuote, sep string, maxLength int) []string {
	words := []string{}
	curWords := []string{}
	curLength := 0
	for _, word := range slice {
		word = strings.TrimSpace(word)
		if len(word) == 0 {
			continue
		}
		if maxLength > 0 && curLength+len(begQuote+word+endQuote+sep) > maxLength {
			words = append(words, strings.Join(curWords, sep))
			curWords = []string{}
			curLength = 0
		} else {
			curWords = append(curWords, begQuote+word+endQuote)
			curLength += len(begQuote + word + endQuote + sep)
		}
	}
	if len(curWords) > 0 {
		words = append(words, strings.Join(curWords, sep))
	}
	return words
}

func JoinQuote(slice []string, begQuote, endQuote, sep string) string {
	words := []string{}
	for _, word := range slice {
		words = append(words, begQuote+word+endQuote)
	}
	return strings.Join(words, sep)
}
