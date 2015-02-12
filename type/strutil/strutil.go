package strutil

import (
	"regexp"
	"strings"
	"unicode"
)

type StrUtil struct {
	RxSpaceBeg   *regexp.Regexp
	RxSpaceEnd   *regexp.Regexp
	RxSpacePunct *regexp.Regexp
	RxDash       *regexp.Regexp
}

func NewStrUtil() StrUtil {
	str := StrUtil{}
	str.RxSpaceBeg = regexp.MustCompile(`^[\s\t\r\n\v\f]*`)
	str.RxSpaceEnd = regexp.MustCompile(`[\s\t\r\n\v\f]*$`)
	str.RxSpacePunct = regexp.MustCompile(`[[:punct:]\s\t\r\n\v\f]`)
	str.RxDash = regexp.MustCompile(`-+`)
	return str
}

func (str *StrUtil) Trim(bytes []byte) []byte {
	bytes = str.RxSpaceBeg.ReplaceAll(bytes, []byte{})
	bytes = str.RxSpaceEnd.ReplaceAll(bytes, []byte{})
	return bytes
}

func Capitalize(s1 string) string {
	s2 := strings.ToLower(s1)
	return ToUpperFirst(s2)
}

func ToLowerFirst(s1 string) string {
	a1 := []rune(s1)
	a1[0] = unicode.ToLower(a1[0])
	return string(a1)
}

func ToUpperFirst(s1 string) string {
	a1 := []rune(s1)
	a1[0] = unicode.ToUpper(a1[0])
	return string(a1)
}
