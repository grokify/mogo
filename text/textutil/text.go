package textutil

import (
	"bufio"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/grokify/mogo/errors/errorsutil"
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

const (
	LineSeparatorN  = "\n"
	LineSeparatorR  = "\r"
	LineSeparatorRN = "\r\n"
)

func ReplaceLineSeparatorFile(infile, outfile, sep string) error {
	r, err := os.OpenFile(infile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return errorsutil.Wrap(err, "open file error")
	}
	defer r.Close()
	w, err := os.Create(outfile)
	if err != nil {
		return errorsutil.Wrap(err, "open file error")
	}
	defer w.Close()
	return ReplaceLineSeparator(r, w, sep)
}

func ReplaceLineSeparator(r io.Reader, w io.Writer, sep string) error {
	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		if _, err := w.Write([]byte(sc.Text() + sep)); err != nil {
			return err
		}
	}
	if err := sc.Err(); err != nil {
		return errorsutil.Wrap(err, "scan reader  error")
	}
	return nil
}
