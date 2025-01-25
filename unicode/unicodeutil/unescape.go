package unicodeutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/grokify/mogo/errors/errorsutil"
)

// Unescape wraps the `strconv.Unquote()` function to provide a method of converting
// \u escaped unicode literals.
func Unescape(s string) (string, error) {
	if !strings.Contains(s, QuotationMark) {
		return strconv.Unquote(QuotationMark + s + QuotationMark)
	} else if !strings.Contains(s, Apostrophe) {
		return strconv.Unquote(Apostrophe + s + Apostrophe)
	} else {
		return s, fmt.Errorf("cannot unescape string (%s)", s)
	}
}

var rxUnicode = regexp.MustCompile(`\\u[0-9A-Za-z]{4}`)

func UnescapeEach(s string) (string, error) {
	var errorsOn []string
	fn := func(s string) string {
		if o, err := Unescape(s); err != nil {
			errorsOn = append(errorsOn, errorsutil.Wrapf(err, "unescape string (%s)", s).Error())
			return s
		} else {
			return o
		}
	}
	out := rxUnicode.ReplaceAllStringFunc(s, fn)
	if len(errorsOn) > 0 {
		return out, fmt.Errorf("errors on [%s]", strings.Join(errorsOn, ", "))
	} else {
		return out, nil
	}
}
