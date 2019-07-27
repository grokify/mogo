// base36 supports Base36 encoding and decoding.
package base36

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/grokify/base36"
	"github.com/grokify/gotilla/math/bigutil"
)

// Encode36 returns an encoded string given a byte array.
func Encode36(ba []byte) string {
	return Encode36HexString(hex.EncodeToString(ba))
}

// Encode36String returns an encoded string given a string.
func Encode36String(s string) string {
	return Encode36([]byte(s))
}

// Encode36HexString returns an encoded string given a string.
func Encode36HexString(s16 string) string {
	i := bigutil.HexToInt(s16)
	return strings.ToLower(base36.EncodeBigInt(i))
}

// Decode36 returns a decoded byte array given an encoded byte array.
func Decode36(b36 []byte) ([]byte, error) {
	return Decode36String(string(b36))
}

// Decode36String returns a decoded byte array given an encoded string.
func Decode36String(s36 string) ([]byte, error) {
	bi := base36.DecodeBigInt(s36)
	return hex.DecodeString(bigutil.IntToHex(bi))
}

// Md5Base36 returns a Base36 encoded MD5 hash of a string.
func Md5Base36(s string) string {
	return fmt.Sprintf("%025s",
		Encode36HexString(fmt.Sprintf("%x", md5.Sum([]byte(s)))))
}

/*
GMP Versions

import "github.com/grokify/gmp"

// Encode36String returns an encoded string given a byte array.
func Encode36(ba []byte) string {
	s16 := hex.EncodeToString(ba)
	bi := gmp.NewInt(0)
	bi.SetString(s16, 16)
	return bi.InBase(36)
}

// Decode36String returns a decoded byte array given an encoded string.
func Decode36String(s36 string) ([]byte, error) {
	bi := gmp.NewInt(0)
	bi.SetString(s36, 36)
	s16 := bi.InBase(16)
	return hex.DecodeString(s16)
}

// Md5Base36 returns a Base36 encoded MD5 hash of a string.
func Md5Base36(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	s16 := fmt.Sprintf("%x", h.Sum(nil))
	bi := gmp.NewInt(0)
	bi.SetString(s16, 16)
	return fmt.Sprintf("%025s", bi.InBase(36))
}
*/
