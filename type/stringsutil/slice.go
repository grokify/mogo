package stringsutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// SliceTrimSpace removes leading and trailing spaces per
// string and also removes empty strings.
func SliceTrimSpace(slice []string) []string {
	trimmed := []string{}
	for _, part := range slice {
		part := strings.TrimSpace(part)
		if len(part) > 0 {
			trimmed = append(trimmed, part)
		}
	}
	return trimmed
}

func SliceIndexOrEmpty(s []string, index uint64) string {
	if int(index) >= len(s) {
		return ""
	}
	return s[index]
}

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
	return strings.Join(SliceTrimSpace(slice), sep)
}

func SliceCondenseRegexps(texts []string, regexps []*regexp.Regexp, replacement string) []string {
	parts := []string{}
	for _, part := range texts {
		for _, rx := range regexps {
			part = rx.ReplaceAllString(part, replacement)
		}
		part = strings.TrimSpace(part)
		if len(part) > 0 {
			parts = append(parts, part)
		}
	}
	return parts
}

func SliceCondensePunctuation(texts []string) []string {
	parts := []string{}
	for _, part := range texts {
		part = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(part, " ")
		part = regexp.MustCompile(`\s+`).ReplaceAllString(part, " ")
		part = strings.TrimSpace(part)
		if len(part) > 0 {
			parts = append(parts, part)
		}
	}
	return parts
}

func SliceCondenseAndQuoteSpace(items []string, quoteLeft, quoteRight string) []string {
	return SliceCondenseAndQuote(items, " ", " ", quoteLeft, quoteRight)
}

func SliceCondenseAndQuote(items []string, trimLeft, trimRight, quoteLeft, quoteRight string) []string {
	newItems := []string{}
	for _, item := range items {
		item = strings.TrimLeft(item, trimLeft)
		item = strings.TrimRight(item, trimRight)
		if len(item) > 0 {
			item = quoteLeft + item + quoteRight
			newItems = append(newItems, item)
		}
	}
	return newItems
}

// SplitTrimSpace splits a string and trims spaces on
// remaining elements.
func SplitTrimSpace(s, sep string) []string {
	split := strings.Split(s, sep)
	strs := []string{}
	for _, str := range split {
		strs = append(strs, strings.TrimSpace(str))
	}
	return strs
}

// SliceToSingleIntOrNeg converts a single element slice
// with a string to an integer or `-1`
func SliceToSingleIntOrNeg(vals []string) int {
	if len(vals) != 1 {
		return -1
	}
	num, err := strconv.Atoi(vals[0])
	if err != nil {
		return -1
	}
	return num
}

// Dedupe returns a string slice with duplicate values
// removed. First observance is kept.
func Dedupe(vals []string) []string {
	deduped := []string{}
	seen := map[string]int{}

	for _, val := range vals {
		if _, ok := seen[val]; ok {
			continue
		}
		seen[val] = 1
		deduped = append(deduped, val)
	}
	return deduped
}

// SliceIndexOf returns the index of an element in a
// string slice. Returns -1 if not found.
func SliceIndexOf(needle string, haystack []string) int {
	for k, v := range haystack {
		if v == needle {
			return k
		}
	}
	return -1 //not found.
}

// SliceIndexOfLcTrimSpace returns the index of an element in a
// string slice. Returns -1 if not found.
func SliceIndexOfLcTrimSpace(needle string, haystack []string) int {
	needle = strings.ToLower(strings.TrimSpace(needle))
	for k, v := range haystack {
		v = strings.ToLower(strings.TrimSpace(v))
		if v == needle {
			return k
		}
	}
	return -1 //not found.
}

func SliceChooseOnePreferredLowerTrimSpace(options, preferenceOrder []string) string {
	if len(options) == 0 {
		return ""
	} else if len(preferenceOrder) == 0 {
		return strings.ToLower(strings.TrimSpace(options[0]))
	}
	optMap := map[string]int{}
	for _, opt := range options {
		opt = strings.ToLower(strings.TrimSpace(opt))
		if len(opt) > 0 {
			optMap[opt] = 1
		}
	}
	for _, pref := range preferenceOrder {
		pref = strings.ToLower(strings.TrimSpace(pref))
		if _, ok := optMap[pref]; ok {
			return pref
		}
	}
	return strings.ToLower(strings.TrimSpace(options[0]))
}

func SliceJoinQuotedMaxLength(slice []string, begQuote, endQuote, sep string, maxLength int) []string {
	words := []string{}
	curWords := []string{}
	curLength := 0
	for _, word := range slice {
		if curLength+len(begQuote+word+endQuote+sep) > maxLength {
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

type JoinCustomConfig struct {
	QuoteBegin string
	QuoteEnd   string
	Separator  string
	MaxLength  int
	TrimSpace  bool
	SkipEmpty  bool
}

func JoinCustom(slice []string, cfg JoinCustomConfig) []string {
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

func SliceJoinQuotedMaxLengthTrimSpaceSkipEmpty(slice []string, begQuote, endQuote, sep string, maxLength int) []string {
	words := []string{}
	curWords := []string{}
	curLength := 0
	for _, word := range slice {
		if curLength+len(begQuote+word+endQuote+sep) > maxLength {
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

func SliceJoinQuoted(slice []string, begQuote, endQuote, sep string) string {
	words := []string{}
	for _, word := range slice {
		words = append(words, begQuote+word+endQuote)
	}
	return strings.Join(words, sep)
}
