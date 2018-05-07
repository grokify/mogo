package stringsutil

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	StringToLower     = "StringToLower"
	SpaceToHyphen     = "SpaceToHyphen"
	SpaceToUnderscore = "SpaceToUnderscore"
)

var rxSpaces = regexp.MustCompile(`\s+`)

// PadLeft prepends a string to a base string until the string
// length is greater or equal to the desired length.
func PadLeft(str string, pad string, length int) string {
	for {
		str = pad + str
		if len(str) >= length {
			return str[0:length]
		}
	}
}

// PadRight appends a string to a base string until the string
// length is greater or equal to the desired length.
func PadRight(str string, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}

// ToLowerFirst lower cases the first letter in the string
func ToLowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

// ToUpperFirst upper cases the first letter in the string
func ToUpperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

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

// SplitCondenseSpace splits a string and trims spaces on
// remaining elements, removing empty elements.
func SplitCondenseSpace(s, sep string) []string {
	split := strings.Split(s, sep)
	strs := []string{}
	for _, str := range split {
		str = strings.TrimSpace(str)
		if len(str) > 0 {
			strs = append(strs, str)
		}
	}
	return strs
}

// CondenseString trims whitespace at the ends of the string
// as well as in between.
func CondenseString(content string, join_lines bool) string {
	if join_lines {
		content = regexp.MustCompile(`\n`).ReplaceAllString(content, " ")
	}
	// Beginning
	content = regexp.MustCompile(`^\s+`).ReplaceAllString(content, "")
	// End
	content = regexp.MustCompile(`\s+$`).ReplaceAllString(content, "")
	// Middle
	content = regexp.MustCompile(`\n[\s\t\r]*\n`).ReplaceAllString(content, "\n")
	// Indentation
	content = regexp.MustCompile(`\n[\s\t\r]*`).ReplaceAllString(content, "\n")
	// Collapse
	content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
	return strings.TrimSpace(content)
}

// TrimSentenceLength trims a string by a max length at word boundaries.
func TrimSentenceLength(sentenceInput string, maxLength int) string {
	if len(sentenceInput) <= maxLength {
		return sentenceInput
	}
	sentenceLen := string(sentenceInput[0:maxLength]) // first350 := string(s[0:350])
	rx_end := regexp.MustCompile(`[[:punct:]][^[[:punct:]]]*$`)
	sentencePunct := rx_end.ReplaceAllString(sentenceLen, "")
	if len(sentencePunct) >= 2 {
		return sentencePunct
	}
	return sentenceLen
}

func JoinTrimSpace(strs []string) string {
	return rxSpaces.ReplaceAllString(strings.Join(strs, " "), " ")
}

// JoinInterface joins an interface and returns a string. It takes
// a join separator, boolean to replace the join separator in the
// string parts and a separator alternate. `stripEmbeddedSep` strips
// separator string found within parts. `stripRepeatedSep` strips
// repeating separators. This flexibility is designed to support
// joining data for both CSVs and paths.
func JoinInterface(arr []interface{}, sep string, stripRepeatedSep bool, stripEmbeddedSep bool, altSep string) string {
	parts := []string{}
	rx := regexp.MustCompile(sep)
	for _, el := range arr {
		part := fmt.Sprintf("%v", el)
		if stripEmbeddedSep {
			part = rx.ReplaceAllString(part, altSep)
		}
		parts = append(parts, part)
	}
	joined := strings.Join(parts, sep)
	if stripRepeatedSep {
		joined = regexp.MustCompile(fmt.Sprintf("%s+", sep)).
			ReplaceAllString(joined, sep)
	}
	return joined
}

func JoinLiterary(slice []string, sep, joinWord string) string {
	switch len(slice) {
	case 0:
		return ""
	case 1:
		return slice[0]
	case 2:
		return slice[0] + " " + joinWord + " " + slice[1]
	default:
		last, rest := slice[len(slice)-1], slice[:len(slice)-1]
		rest = append(rest, joinWord+" "+last)
		return strings.Join(rest, sep+" ")
	}
}

func JoinLiteraryQuote(slice []string, leftQuote, rightQuote, sep, joinWord string) string {
	newSlice := SliceCondenseAndQuoteSpace(slice, leftQuote, rightQuote)
	switch len(newSlice) {
	case 0:
		return ""
	case 1:
		return newSlice[0]
	case 2:
		return newSlice[0] + " " + joinWord + " " + newSlice[1]
	default:
		last, rest := newSlice[len(newSlice)-1], newSlice[:len(newSlice)-1]
		rest = append(rest, joinWord+" "+last)
		return strings.Join(rest, sep+" ")
	}
}

func FormatString(s string, options []string) string {
	for _, opt := range options {
		switch strings.TrimSpace(opt) {
		case StringToLower:
			s = strings.ToLower(s)
		case SpaceToHyphen:
			s = regexp.MustCompile(`[\s-]+`).ReplaceAllString(s, "-")
		case SpaceToUnderscore:
			s = regexp.MustCompile(`[\s_]+`).ReplaceAllString(s, "_")
		}
	}
	return s
}
