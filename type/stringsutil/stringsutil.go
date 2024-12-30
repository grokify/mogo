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
	// lowerAZ           = "abcdefghijklmnopqrstuvwxyz"
	// upperAZ           = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerUpper = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	upperLower = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

type Stringable interface {
	String() string
}

type StringableWithErr interface {
	String() (string, error)
}

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

func IsLower(s string) bool {
	return s == strings.ToLower(s)
}

func IsUpper(s string) bool {
	return s == strings.ToUpper(s)
}

// Capitalize returns a string with the first character
// capitalized and the rest lower cased.
func Capitalize(s1 string) string {
	return ToUpperFirst(s1, true)
}

// ToLowerFirst lower cases the first letter in the string
func ToLowerFirst(s1 string) string {
	a1 := []rune(s1)
	a1[0] = unicode.ToLower(a1[0])
	return string(a1)
}

// ToUpperFirst upper cases the first letter in the string
func ToUpperFirst(s1 string, lowerRest bool) string {
	if lowerRest {
		s1 = strings.ToLower(s1)
	}
	a1 := []rune(s1)
	if len(a1) > 0 {
		a1[0] = unicode.ToUpper(a1[0])
	}
	return string(a1)
}

// ToBool converts a string to a boolean value
// converting "f", "false", "0" and the empty string
// to false with everything else being true.
func ToBool(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" || v == "0" || v == "f" || v == "false" {
		return false
	}
	return true
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

var rxSpace = regexp.MustCompile(`\s+`)

// RemoveSpaces eliminates all spaces in a string.
func RemoveSpaces(input string) string {
	return rxSpace.ReplaceAllString(input, "")
}

// SplitCondenseSpace splits a string and trims spaces on
// remaining elements, removing empty elements.
func SplitCondenseSpace(s, sep string) []string {
	split := strings.Split(s, sep)
	strs := []string{}
	for _, str := range split {
		if str = strings.TrimSpace(str); len(str) > 0 {
			strs = append(strs, str)
		}
	}
	return strs
}

// CondenseString trims whitespace at the ends of the string
// as well as in between.
func CondenseString(content string, joinLines bool) string {
	if joinLines {
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

// CondenseSpace removes extra spaces.
func CondenseSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func StripControl(s string) string { return rxControl.ReplaceAllString(s, "") }

func StripSubstring(s, substr string, insensitive bool) string {
	var rx *regexp.Regexp
	if insensitive {
		rx = regexp.MustCompile(`(?i)` + regexp.QuoteMeta(substr))
	} else {
		rx = regexp.MustCompile(regexp.QuoteMeta(substr))
	}
	return rx.ReplaceAllString(s, "")
}

/*
func OrDefault(s, defaultValue string) string {
	if len(s) == 0 {
		return defaultValue
	}
	return s
}
*/

func FirstNonEmpty(vals ...string) string {
	for _, val := range vals {
		if len(val) > 0 {
			return val
		}
	}
	return ""
}

// TrimSpaceOrDefault trims spaces and replaces default value if
// result is empty string.
func TrimSpaceOrDefault(str, defaultValue string) string {
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		str = strings.TrimSpace(defaultValue)
	}
	return str
}

// TrimSentenceLength trims a string by a max length at word boundaries.
func TrimSentenceLength(sentenceInput string, maxLength int) string {
	if len(sentenceInput) <= maxLength {
		return sentenceInput
	}
	sentenceLen := sentenceInput[0:maxLength] // first350 := string(s[0:350])
	rxEnd := regexp.MustCompile(`[[:punct:]][^[[:punct:]]]*$`)
	sentencePunct := rxEnd.ReplaceAllString(sentenceLen, "")
	if len(sentencePunct) >= 2 {
		return sentencePunct
	}
	return sentenceLen
}

// FirstNotEmptyTrimSpace returns the first non-empty string
// after applying `strings.TrimSpace()`.`
func FirstNotEmptyTrimSpace(candidates ...string) string {
	for _, s := range candidates {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			return s
		}
	}
	return ""
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
func JoinInterface(arr []any, sep string, stripRepeatedSep bool, stripEmbeddedSep bool, altSep string) string {
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
			newParts = append(newParts, ToUpperFirst(p, true))
		}
	}
	return strings.Join(newParts, "")
}

func ToOpposite(s string) string {
	return xstrings.Translate(s, lowerUpper, upperLower)
}

var (
	rxRN = regexp.MustCompile(`\r\n`)
	rxR  = regexp.MustCompile(`\r`)
)

func NewlineToLinux(input string) string {
	return rxR.ReplaceAllString(rxRN.ReplaceAllString(input, "\n"), "\n")
}

// EmptyError takes a string and error, returning
// the string value or an empty string if an error
// is encountered. It is used for simplifying code
// that returns a value or an error if not present.
func EmptyError(s string, err error) string {
	if err != nil {
		return ""
	}
	return s
}

// RemoveNonPrintable removes non-printable characters from a string.
func RemoveNonPrintable(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		} else {
			return -1
		}
	}, s)
}

// StripChars removes chars specified by `cutset` while maintaining order of remaining
// chars and shortening string per removed chars.
func StripChars(s, cutset string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(cutset, r) {
			return -1
		} else {
			return r
		}
	}, s)
}

// UniqueRunes checks to see if a string's runes are unique.
func UniqueRunes(s string) bool {
	v := map[rune]bool{}
	for _, r := range s {
		if v[r] {
			return false
		} else {
			v[r] = true
		}
	}
	return len(s) == len(v)
}

// Repeat returns atring of length `length` by repeating string `s`. If `length` is less than
// then length of `s`, the result is cut to `length`.
func Repeat(s string, length uint) string {
	if length == 0 {
		return ""
	}
	str := ""
	for {
		str += s
		l := uint(len(str))
		if l == length {
			return str
		} else if l > length {
			return str[:int(length)]
		}
	}
}
