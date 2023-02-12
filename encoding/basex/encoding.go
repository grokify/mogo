package basex

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	AlphabetBase10          = "0123456789"
	AlphabetBase16          = "0123456789abcdef"
	AlphabetBase26          = "abcdefghijklmnopqrstuvwxyz"
	AlphabetBase31          = "0123456789BCDFGHJKLMNPQRSTVWXYZ"  // Gaming, no foul words
	AlphabetBase32          = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567" // RFC-4648: https://en.wikipedia.org/wiki/Base32
	AlphabetBase32Geohash   = "0123456789bcdefghjkmnpqrstuvwxyz"
	AlphabetBase32Hex       = "0123456789ABCDEFGHIJKLMNOPQRSTUV"
	AlphabetBase32Wordsafe  = "23456789CFGHJMPQRVWXcfghjmpqrvwx" // Word-safe alphabet to avoid forming words
	AlphabetBase32Z         = "ybndrfg8ejkmcpqxot1uwisza345h769"
	AlphabetBase36          = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphabetBase45          = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:"            // https://datatracker.ietf.org/doc/rfc9285/
	AlphabetBase56          = "0123456789ABCEFGHJKLMNPRSTUVWXYZabcdefghjklmnpqrstuvwxyz" // See: https://github.com/tep/encoding-base56
	AlphabetBase56Alt       = "23456789abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ" // used by Java and PHP
	AlphabetBase56Java      = AlphabetBase56Alt                                          // See: https://github.com/tep/encoding-base56
	AlphabetBase56PHP       = AlphabetBase56Alt                                          // See: https://github.com/tep/encoding-base56
	AlphabetBase56Python3   = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"
	AlphabetBase58          = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz" // Bitcoin alphabet: https://en.bitcoin.it/wiki/Base58Check_encoding
	AlphabetBase58Flickr    = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	AlphabetBase58GMP       = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuv"
	AlphabetBase58Ripple    = "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"
	AlphabetBase62          = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" // ASCII table, used by GMP
	AlphabetBase62Gobigint  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphabetBase62LUN       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	AlphabetBase62ULN       = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	AlphabetBase63          = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
	AlphabetBase63Roblox    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_" // Using traditional Roblox ordering (not for encoding)
	AlphabetBase64          = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	AlphabetBase64IMAP      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+,"
	AlphabetBase64URL       = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	AlphabetBase70AWSS3Safe = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!-_.*'()"                // https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-keys.html
	AlphabetBase85          = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!#$%&()*+-;<=>?@^_`{|}~" // RFC-1924 https://en.wikipedia.org/wiki/Ascii85
)

func AlphabetDefault(base uint) string {
	switch base {
	case 16:
		return AlphabetBase16
	case 31:
		return AlphabetBase31
	case 32:
		return AlphabetBase32
	case 36:
		return AlphabetBase36
	case 45:
		return AlphabetBase45
	case 56:
		return AlphabetBase56
	case 58:
		return AlphabetBase58
	case 62:
		return AlphabetBase62
	case 63:
		return AlphabetBase63
	case 64:
		return AlphabetBase64
	case 85:
		return AlphabetBase85
	default:
		if base < 85 {
			return AlphabetBase85[:base]
		} else {
			panic("alphabet desired larger than max 85")
		}
	}
}

var (
	ErrInvalidAlphaNumericAlphabet  = errors.New("invalid alphanumeric alphabet")
	ErrInvalidAlphabetHasDuplicates = errors.New("invalid alphabet has duplicates")
	ErrInvalidAlphabetIsEmpty       = errors.New("invalid alphabet is empty")
)

// RuneMaps produces maps of alphabets.
func RuneMaps(alphabet string) (map[int]rune, map[rune]int) {
	m := map[int]rune{}
	n := map[rune]int{}
	for i, l := range alphabet {
		m[i] = l
		n[l] = i
	}
	return m, n
}

// AlphabetShuffled shuffles an alphabet to provide a random ordering.
func AlphabetShuffled(alphabet string) string {
	letters := strings.Split(alphabet, "")
	rand.Shuffle(len(letters), func(i, j int) {
		letters[i], letters[j] = letters[j], letters[i]
	})
	return strings.Join(letters, "")
}

func IsAlphaNumericAlphabet(alphabet string) bool {
	if !stringsutil.IsAlphaNumeric(alphabet) {
		return false
	} else if !slicesutil.UniqueValues(strings.Split(alphabet, "")) {
		return false
	}
	return true
}

// ValidAlphabet checks to see if string `s` is within the supplied alphabet.
func ValidAlphabet(alphabet, s string) bool {
	_, n := RuneMaps(alphabet)
	for _, l := range s {
		if _, ok := n[l]; !ok {
			return false
		}
	}
	return true
}

// ValidAlphabetMap checks to see if string `s` is within the supplied alphabet.
// Prefer this over `ValidAlphabet()` when calling many times.
func ValidAlphabetMap(alphabet map[rune]int, s string) bool {
	for _, l := range s {
		if _, ok := alphabet[l]; !ok {
			return false
		}
	}
	return true
}

// ValidAlphabets validates a pair of alphabets with the second being optional.
func ValidAlphabets(a1, a2 string, a2optional bool) error {
	if a1 == "" {
		return errors.New("alphabet 1 cannot be empty")
	} else if !stringsutil.UniqueRunes(a1) {
		return errors.New("alphabet 1 chars are not unique")
	} else if a2optional && a2 == "" {
		return nil
	} else if len(a1) != len(a2) {
		return errors.New("alphabet 1 and alphabet 2 length mismatch")
	} else if !stringsutil.UniqueRunes(a2) {
		return errors.New("alphabet 2 chars are not unique")
	}
	return nil
}
