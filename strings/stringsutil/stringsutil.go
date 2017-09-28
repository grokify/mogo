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
