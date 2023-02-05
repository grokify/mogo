package stringsutil

import (
	"regexp"
	"strings"
)

const (
	NewlinePOSIX       = "\n"
	NewlineMSFT        = "\r\n"
	NewlineAPPLClassic = "\r"
)

var (
	rxCarriageReturn         = regexp.MustCompile(`\r`)
	rxCarriageReturnLineFeed = regexp.MustCompile(`\r\n`)
	rxLinefeedsGTEThree      = regexp.MustCompile(`\n\n\n+`)
)

func ToLineFeeds(s string) string {
	return rxCarriageReturn.ReplaceAllString(
		rxCarriageReturnLineFeed.ReplaceAllString(s, "\n"),
		"\n")
}

func SplitLines(s string) []string {
	return strings.Split(ToLineFeeds(s), "\n")
}

// CondenseLines improves text layout by (a) remove staring and trailing spaces from eacn line,
// (b) remove starting and training new lines before and after characters, and (c) ensuring that
// there is a max of 2 consecutive line feeds.
func CondenseLines(s, outNewline string) string {
	if outNewline == "" {
		outNewline = NewlinePOSIX
	}
	l := strings.Split(strings.TrimSpace(CarriageReturnsToLinefeeds(s)), NewlinePOSIX)
	for i, li := range l {
		l[i] = strings.TrimSpace(li)
	}
	s = strings.Join(l, outNewline)
	return rxLinefeedsGTEThree.ReplaceAllString(s, outNewline+outNewline)
}

// CarriageReturnsToLinefeeds replaces `\r\n` with `\n` followed by replacing `\r` by `\n`.
func CarriageReturnsToLinefeeds(s string) string {
	return rxCarriageReturn.ReplaceAllString(
		rxCarriageReturnLineFeed.ReplaceAllString(s, NewlinePOSIX),
		NewlinePOSIX)
}
