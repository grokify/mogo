package base36

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/grokify/gmp"
)

func Encode36(ba []byte) string {
	s16 := hex.EncodeToString(ba)
	bi := gmp.NewInt(0)
	bi.SetString(s16, 16)
	return bi.InBase(36)
}

func Encode36String(s string) string {
	return Encode36([]byte(s))
}

func Decode36(b36 []byte) ([]byte, error) {
	return Decode36String(string(b36))
}

func Decode36String(s36 string) ([]byte, error) {
	bi := gmp.NewInt(0)
	bi.SetString(s36, 36)
	s16 := bi.InBase(16)
	return hex.DecodeString(s16)
}

func Md5Base36(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	s16 := fmt.Sprintf("%x", h.Sum(nil))
	bi := gmp.NewInt(0)
	bi.SetString(s16, 16)
	return fmt.Sprintf("%025s", bi.InBase(36))
}
