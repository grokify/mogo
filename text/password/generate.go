package password

import (
	"strings"

	"github.com/grokify/mogo/crypto/randutil"
	"github.com/grokify/mogo/encoding"
	"github.com/grokify/mogo/type/stringsutil"
)

// Generate creates passwords using options selected in `GenerateOpts`.
func Generate(opts GenerateOpts, l uint) string {
	return randutil.RandString(BuildAlphabet(opts), l)
}

const (
	AlphabetSymbols   = "!@#$%&*?"
	AlphabetAmbiguous = "{}[]()/\\'\"`,;:.<>"
	AlphabetSimilar   = "iI1LoO0"
)

type GenerateOpts struct {
	InclLower     bool
	InclUpper     bool
	InclNumbers   bool
	InclSymbols   bool
	InclAmbiguous bool
	ExclSimilar   bool
}

// BuildAlphabet builds an alphabet that's useful for passwords.
func BuildAlphabet(opts GenerateOpts) string {
	var alphabet string
	if opts.InclLower {
		alphabet += strings.ToLower(encoding.AlphabetBase26)
	}
	if opts.InclNumbers {
		alphabet += encoding.AlphabetBase10
	}
	if opts.InclUpper {
		alphabet += strings.ToUpper(encoding.AlphabetBase26)
	}
	if opts.InclSymbols {
		alphabet += AlphabetSymbols
	}
	if opts.InclAmbiguous {
		alphabet += AlphabetAmbiguous
	}
	if opts.ExclSimilar {
		alphabet = stringsutil.StripChars(alphabet, AlphabetSimilar)
	}
	return alphabet
}
