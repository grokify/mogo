package stringsutil

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/huandu/xstrings"
)

const (
	StringToLower     = "StringToLower"
	SpaceToHyphen     = "SpaceToHyphen"
	SpaceToUnderscore = "SpaceToUnderscore"
	UpperAZ           = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowerAZ           = "abcdefghijklmnopqrstuvwxyz"
	LowerUpper        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	UpperLower        = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

var (
	rxControl = regexp.MustCompile(`[[:cntrl:]]`)
	rxSpaces  = regexp.MustCompile(`\s+`)
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

// Capitalize returns a string with the first character
// capitalized and the rest lower cased.
func Capitalize(s1 string) string {
	s2 := strings.ToLower(s1)
	return ToUpperFirst(s2)
}

// ToLowerFirst lower cases the first letter in the string
func ToLowerFirst(s1 string) string {
	a1 := []rune(s1)
	a1[0] = unicode.ToLower(a1[0])
	return string(a1)
	/*
		if s == "" {
			return ""
		}
		r, n := utf8.DecodeRuneInString(s)
		return string(unicode.ToLower(r)) + s[n:]
	*/
}

// ToUpperFirst upper cases the first letter in the string
func ToUpperFirst(s1 string) string {
	a1 := []rune(s1)
	a1[0] = unicode.ToUpper(a1[0])
	return string(a1)
	/*
			if s == "" {
			return ""
		}
		r, n := utf8.DecodeRuneInString(s)
		return string(unicode.ToUpper(r)) + s[n:]
	*/
}

// ToBool converts a string to a boolean value
// looking for the string "true" in any case.
func ToBool(v string) bool {
	if strings.TrimSpace(strings.ToLower(v)) == "true" {
		return true
	}
	return false
}

func SubstringIsSuffix(s1, s2 string) bool {
	len1 := len(s1)
	len2 := len(s2)
	idx := strings.Index(s1, s2)
	if len1 >= len2 && idx > -1 && idx == (len1-len2) {
		return true
	}
	return false
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

func StripControl(s string) string { return rxControl.ReplaceAllString(s, "") }

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

func JoinStringsTrimSpaceToLowerSort(strs []string, sep string) string {
	wip := []string{}
	for _, s := range strs {
		s = strings.ToLower(strings.TrimSpace(s))
		if len(s) > 0 {
			wip = append(wip, s)
		}
	}
	sort.Strings(wip)
	return strings.Join(wip, sep)
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

// CommonInitialisms is the listed by Go Lint.
const CommonInitialisms = "ACL,API,ASCII,CPU,CSS,DNS,EOF,GUID,HTML,HTTP,HTTPS,ID,IP,JSON,LHS,QPS,RAM,RHS,RPC,SLA,SMTP,SQL,SSH,TCP,TLS,TTL,UDP,UI,UID,UUID,URI,URL,UTF8,VM,XML,XMPP,XSRF,XSS"

// CommonInitialismsMap returns map[string]bool of upper case initialisms.
func CommonInitialismsMap() map[string]bool {
	ciMap := map[string]bool{}
	commonInitialisms := strings.Split(CommonInitialisms, ",")
	for _, ci := range commonInitialisms {
		ciMap[ci] = true
	}
	return ciMap
}

// StringToConstant is used to generate constant names for code generation.
// It uses the commonInitialisms in Go Lint.
func StringToConstant(s string) string {
	newParts := []string{}
	parts := strings.Split(s, "_")
	ciMap := CommonInitialismsMap()
	for _, p := range parts {
		pUp := strings.ToUpper(p)
		if _, ok := ciMap[pUp]; ok {
			newParts = append(newParts, pUp)
		} else {
			newParts = append(newParts, ToUpperFirst(strings.ToLower(p)))
		}
	}
	return strings.Join(newParts, "")
}

func ToOpposite(s string) string {
	return xstrings.Translate(s, LowerUpper, UpperLower)
}
