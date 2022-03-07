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
	GUIDPattern = `^([0-9a-fA-F]{8})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{12})$`
	guidReplace = "${1}-${2}-${3}-${4}-${5}"
)

var (
	rxGUID   = regexp.MustCompile(GUIDPattern)
	rxHyphen = regexp.MustCompile(`-`)
)

// ValidGUIDHex checks to see if a string is a valid GUID.
func ValidGUIDHex(guid string) bool {
	return rxGUID.MatchString(guid)
}

// GUIDToBigInt converts a GUID string, with or with out hypens, to a *big.Int.
func GUIDToBigInt(guid string) (*big.Int, error) {
	if !ValidGUIDHex(guid) {
		return nil, fmt.Errorf("not a valid GUID [%s]", guid)
	}
	bi := big.NewInt(0)
	bi.SetString(rxHyphen.ReplaceAllString(guid, ""), 16)
	return bi, nil
}

// GUIDToBase58 converts a GUID string to a Base58 string using the Bitcoin alphabet.
func GUIDToBase58(guid string) (string, error) {
	bi, err := GUIDToBigInt(guid)
	if err != nil {
		return "", err
	}
	return string(bitcoinmath.Big2Base58(bi)), nil
}

// Base58ToGUID converts a Base58 string to a GUID string, with or without hyphens, using the Bitcoin alphabet.
func Base58ToGUID(b58str string, inclHyphen bool) (string, error) {
	b58 := bitcoinmath.Base58(b58str)
	bi := b58.Base582Big()

	guid := fmt.Sprintf("%032s", bi.Text(16))

	if len(guid) != 32 {
		return "", fmt.Errorf("error converting base58 string to hex: %v", b58str)
	}

	if inclHyphen {
		guid = rxGUID.ReplaceAllString(guid, guidReplace)
	}

	return guid, nil
}
