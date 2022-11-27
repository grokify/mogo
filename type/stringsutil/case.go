package stringsutil

import "golang.org/x/text/cases"

var defaultCaser = cases.Fold()

// EqualFoldFull provides "full Unicode case-folding", unlike `strings.EqualFold` which
// provides "simple Unicode case-folding". If `caser` is set to `nil`, the default caser
// with no additional `cases.Option` is used.
func EqualFoldFull(s, t string, caser *cases.Caser) bool {
	if caser != nil {
		return caser.String(s) == caser.String(t)
	}
	return defaultCaser.String(s) == defaultCaser.String(t)
}
