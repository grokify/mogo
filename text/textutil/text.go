package textutil

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func RemoveDiacritics(s string) (string, error) {
	// Should å -> aa: https://stackoverflow.com/questions/11248467/convert-unicode-to-double-ascii-letters-in-python-%C3%9F-ss
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, strings.Replace(s, "\u00df", "ss", -1))
	return result, err
}
