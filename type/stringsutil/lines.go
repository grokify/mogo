package stringsutil

import (
	"regexp"
	"strings"
)

var (
	rxCarriageReturn         = regexp.MustCompile(`\r`)
	rxCarriageReturnLineFeed = regexp.MustCompile(`\r\n`)
)

func ToLineFeeds(s string) string {
	return rxCarriageReturn.ReplaceAllString(
		rxCarriageReturnLineFeed.ReplaceAllString(s, "\n"),
		"\n")
}

func SplitLines(s string) []string {
	return strings.Split(ToLineFeeds(s), "\n")
}
