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

var rxMultiLinefeeds = regexp.MustCompile(`\n\n\n+`)

// CondenseLines improves text layout by (a) remove staring and trailing spaces from eacn line,
// (b) remove starting and training new lines before and after characters, and (c) ensuring that
// there is a max of 2 consecutive line feeds.
func CondenseLines(s string) string {
	s = strings.TrimSpace(s)
	l := strings.Split(s, "\n")
	for i, li := range l {
		l[i] = strings.TrimSpace(li)
	}
	s = strings.Join(l, "\n")
	return rxMultiLinefeeds.ReplaceAllString(s, "\n\n")
}
