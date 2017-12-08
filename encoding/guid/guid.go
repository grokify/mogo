// encoding/guid supports encoding and decoding Guid values.
package guid

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/grokify/bitcoinmath"
)

const (
	// GuidPattern is a regexp pattern for GUIDs.
	GuidPattern = `^([0-9a-fA-F]{8})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{12})$`
	guidReplace = "${1}-${2}-${3}-${4}-${5}"
)

var (
	rxGuid   = regexp.MustCompile(GuidPattern)
	rxHyphen = regexp.MustCompile(`-`)
)

// ValidGuidHex checks to see if a string is a valid GUID.
func ValidGuidHex(guid string) bool {
	return rxGuid.MatchString(guid)
}

// GuidToBigInt converts a GUID string, with or with out hypens, to a *big.Int.
func GuidToBigInt(guid string) (*big.Int, error) {
	if !ValidGuidHex(guid) {
		return nil, fmt.Errorf("Not a valid Guid: %v\n", guid)
	}
	bi := big.NewInt(0)
	bi.SetString(rxHyphen.ReplaceAllString(guid, ""), 16)
	return bi, nil
}

// GuidToBase58 converts a GUID string to a Base58 string using the Bitcoin alphabet.
func GuidToBase58(guid string) (string, error) {
	bi, err := GuidToBigInt(guid)
	if err != nil {
		return "", err
	}
	return string(bitcoinmath.Big2Base58(bi)), nil
}

// Base58ToGuid converts a Base58 string to a GUID string, with or without hyphens, using the Bitcoin alphabet.
func Base58ToGuid(b58str string, inclHyphen bool) (string, error) {
	b58 := bitcoinmath.Base58(b58str)
	bi := b58.Base582Big()

	guid := fmt.Sprintf("%032s", bi.Text(16))

	if len(guid) != 32 {
		return "", fmt.Errorf("Error converting base58 string to hex: %v", b58str)
	}

	if inclHyphen {
		guid = rxGuid.ReplaceAllString(guid, guidReplace)
	}

	return guid, nil
}
