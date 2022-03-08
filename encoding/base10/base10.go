// base10 supports Base10 encoding.
package base10

import (
	"crypto/md5" // #nosec G501
	"math/big"
)

// Encode adapted from https://stackoverflow.com/questions/28128285/
func Encode(bytes []byte) *big.Int {
	bi := big.NewInt(0)
	h := md5.New()
	h.Write(bytes)
	bi.SetBytes(h.Sum(nil))
	return bi
}
