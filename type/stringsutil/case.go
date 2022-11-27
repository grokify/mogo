package stringsutil

import "golang.org/x/text/cases"

// EqualFoldFull provides "full Unicode case-folding", unlike `strings.EqualFold` which
// provides "simple Unicode case-folding". This is fine to run with no options.
func EqualFoldFull(s, t string, opts ...cases.Option) bool {
	c := cases.Fold(opts...)
	return c.String(s) == c.String(t)
}
