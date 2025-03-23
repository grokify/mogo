package stringsutil

import (
	"regexp"
	"strings"
)

const (
	newlinePOSIX       = "\n"
	newlineMSFT        = "\r\n"
	newlineAPPLClassic = "\r"
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

// CondenseLines improves text layout by (a) remove staring and trailing spaces from eacn line,
// (b) remove starting and training new lines before and after characters, and (c) ensuring that
// there is a max of 2 consecutive line feeds.
func CondenseLines(s, outNewline string) string {
	if outNewline == "" {
		outNewline = newlinePOSIX
	}
	l := strings.Split(strings.TrimSpace(CarriageReturnsToLinefeeds(s)), newlinePOSIX)
	for i, li := range l {
		l[i] = strings.TrimSpace(li)
	}
	s = strings.Join(l, outNewline)
	return rxLinefeedsGTEThree.ReplaceAllString(s, outNewline+outNewline)
}

// CarriageReturnsToLinefeeds replaces `\r\n` with `\n` followed by replacing `\r` by `\n`.
func CarriageReturnsToLinefeeds(s string) string {
	return rxCarriageReturn.ReplaceAllString(
		rxCarriageReturnLineFeed.ReplaceAllString(s, newlinePOSIX),
		newlinePOSIX)
}
