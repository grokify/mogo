package stringsutil

import (
	"regexp"
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

// CondenseString trims whitespace at the ends of the string
// as well as in between.
func CondenseString(content string, join_lines bool) string {
	rx_beg := regexp.MustCompile(`^\s+`)
	rx_end := regexp.MustCompile(`\s+$`)
	rx_mid := regexp.MustCompile(`\n[\s\t\r]*\n`)
	rx_pre := regexp.MustCompile(`\n[\s\t\r]*`)
	rx_spc := regexp.MustCompile(`\s+`)
	content = rx_beg.ReplaceAllString(content, "")
	content = rx_end.ReplaceAllString(content, "")
	content = rx_mid.ReplaceAllString(content, "\n")
	content = rx_pre.ReplaceAllString(content, "\n")
	content = rx_spc.ReplaceAllString(content, " ")
	if join_lines {
		rx_join := regexp.MustCompile(`\n`)
		content = rx_join.ReplaceAllString(content, " ")
	}
	return content
}
