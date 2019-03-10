package stringsutil

import (
	"regexp"
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

func InterfaceToSliceString(s interface{}) []string {
	ss := s.([]interface{})
	a := []string{}
	for _, i := range ss {
		a = append(a, i.(string))
	}
	return a
}
