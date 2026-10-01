// Package licenses provides helpers for identifying software license texts.
package licenses

import (
	"regexp"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
)

// mitTerms is the MIT License text after the title and copyright notice.
const mitTerms = `Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.`

// rxMITHeaderLine matches lines allowed before the MIT terms: a title such
// as "MIT License" or "The MIT License (MIT)", or a copyright notice.
var rxMITHeaderLine = regexp.MustCompile(`(?i)^(?:(?:the\s+)?mit\s+license(?:\s+\(mit\))?|(?:copyright|\(c\)|©).*)$`)

// IsMIT reports whether s is the MIT License: the standard MIT terms,
// optionally preceded by a title and copyright lines, with nothing after
// them. Whitespace, line wrapping, and letter case are ignored, so a
// reflowed LICENSE file matches, but any added or altered terms do not.
func IsMIT(s string) bool {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if equalNormalized(strings.Join(lines[i:], "\n"), mitTerms) {
			return true
		}
		if !rxMITHeaderLine.MatchString(line) {
			return false
		}
	}
	return false
}

// equalNormalized compares s and t ignoring whitespace differences and case.
func equalNormalized(s, t string) bool {
	return stringsutil.EqualFoldFull(
		strings.Join(strings.Fields(s), " "),
		strings.Join(strings.Fields(t), " "),
		nil)
}
