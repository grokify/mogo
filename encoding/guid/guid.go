// encoding/guid supports encoding and decoding Guid values.
package guid

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/grokify/bitcoinmath"
)

const (
	GuidPattern = `^([0-9a-fA-F]{8})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{12})$`
	GuidReplace = "${1}-${2}-${3}-${4}-${5}"
)

var (
	rxGuid   = regexp.MustCompile(GuidPattern)
	rxHyphen = regexp.MustCompile(`-`)
)

func ValidGuidHex(guid string) bool {
	return rxGuid.MatchString(guid)
}

func GuidToBigInt(guid string) (*big.Int, error) {
	if !ValidGuidHex(guid) {
		return nil, fmt.Errorf("Not a valid Guid: %v\n", guid)
	}
	hexstr := rxHyphen.ReplaceAllString(guid, "")
	bi := big.NewInt(0)
	bi.SetString(hexstr, 16)
	return bi, nil
}

func GuidToBase58(guid string) (string, error) {
	bi, err := GuidToBigInt(guid)
	if err != nil {
		return "", err
	}
	return string(bitcoinmath.Big2Base58(bi)), nil
}

func Base58ToGuid(b58str string, inclHyphen bool) (string, error) {
	b58 := bitcoinmath.Base58(b58str)
	bi := b58.Base582Big()

	guid := fmt.Sprintf("%032s", bi.Text(16))

	if len(guid) != 32 {
		return "", fmt.Errorf("Error converting base58 string to hex: %v", b58str)
	}

	if inclHyphen {
		guid = rxGuid.ReplaceAllString(guid, GuidReplace)
	}

	return guid, nil
}
