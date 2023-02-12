package password

import (
	"strings"

	"github.com/grokify/mogo/crypto/randutil"
	"github.com/grokify/mogo/encoding/basex"
	"github.com/grokify/mogo/type/stringsutil"
)

// Generate creates passwords using options selected in `GenerateOpts`.
func Generate(opts GenerateOpts) string {
	return randutil.RandString(opts.Alphabet(), opts.Length)
}

const (
	AlphabetSymbols   = "!@#$%&*?"
	AlphabetAmbiguous = "{}[]()/\\'\"`,;:.<>"
	AlphabetSimilar   = "iI1LoO0"
)

type GenerateOpts struct {
	Length        uint
	InclLower     bool
	InclUpper     bool
	InclNumbers   bool
	InclSymbols   bool
	InclAmbiguous bool
	ExclSimilar   bool
}

// Alphabet builds an alphabet that's useful for passwords.
func (opts GenerateOpts) Alphabet() string {
	var alphabet string
	if opts.InclLower {
		alphabet += strings.ToLower(basex.AlphabetBase26)
	}
	if opts.InclNumbers {
		alphabet += basex.AlphabetBase10
	}
	if opts.InclUpper {
		alphabet += strings.ToUpper(basex.AlphabetBase26)
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
