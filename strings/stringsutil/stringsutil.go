package stringsutil

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
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

func JoinInterface(arr []interface{}, sep string) string {
	parts := []string{}
	for _, el := range arr {
		parts = append(parts, fmt.Sprintf("%v", el))
	}
	return strings.Join(parts, sep)
}
