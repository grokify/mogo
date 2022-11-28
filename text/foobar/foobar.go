package foobar

import "strings"

const VarsString = "foo, bar, baz, qux, quux, corge, grault, garply, waldo, fred, plugh, xyzzy, thud"

// Vars retursn a list of foobar metasyntactic variables as seen in this list: https://ascii.jp/elem/000/000/061/61404/
func Vars() []string {
	return strings.Split(VarsString, ", ")
}
