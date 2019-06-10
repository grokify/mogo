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
