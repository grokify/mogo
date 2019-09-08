// md5 supports MD5 hashes in various formats.
package md5

import (
	cryptomd5 "crypto/md5"
	"fmt"
	"math/big"

	"github.com/grokify/gotilla/encoding/base36"
)

// Md5Base36Length is the length for a MD5 Base36 string
const (
	Md5Base36Length int    = 25
	md5Base36Format string = `%025s`
	Md5Base10Length int    = 39
	md5Base10Format string = `%039s`
)

// Md5Base10 returns a Base10 encoded MD5 hash of a string.
func Md5Base10(s string) string {
	i := new(big.Int)
	i.SetString(fmt.Sprintf("%x", cryptomd5.Sum([]byte(s))), 16)
	return fmt.Sprintf(md5Base10Format, i.String())
}

// Md5Base36 returns a Base36 encoded MD5 hash of a string.
func Md5Base36(s string) string {
	return fmt.Sprintf(md5Base36Format,
		base36.Encode36HexString(fmt.Sprintf("%x", cryptomd5.Sum([]byte(s)))))
}
