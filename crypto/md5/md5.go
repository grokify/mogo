// md5 supports MD5 hashes in various formats.
package md5

import (
	cryptomd5 "crypto/md5" // #nosec G501
	"fmt"
	"math/big"

	"github.com/grokify/mogo/type/stringsutil"
)

// Md5Base36Length is the length for a MD5 Base36 string
const (
	// md5Base62Length int    = 22
	md5Base62Format string = `%022s`
	// md5Base36Length int    = 25
	md5Base36Format string = `%025s`
	// md5Base10Length int    = 39
	md5Base10Format string = `%039s`
)

// Md5Base10 returns a Base10 encoded MD5 hash of a string.
func Md5Base10(s string) string {
	i := new(big.Int)
	i.SetString(fmt.Sprintf("%x", cryptomd5.Sum([]byte(s))), 16) // #nosec G401
	return fmt.Sprintf(md5Base10Format, i.String())
}

// Md5Base36 returns a Base36 encoded MD5 hash of a string.
func Md5Base36(s string) string {
	i := new(big.Int)
	i.SetString(fmt.Sprintf("%x", cryptomd5.Sum([]byte(s))), 16) // #nosec G401
	return fmt.Sprintf(md5Base36Format, i.Text(36))
}

// Md5Base62 returns a Base62 encoded MD5 hash of a string.
// This uses the Golang alphabet [0-9a-zA-Z].
func Md5Base62(s string) string {
	i := new(big.Int)
	i.SetString(fmt.Sprintf("%x", cryptomd5.Sum([]byte(s))), 16) // #nosec G401
	return fmt.Sprintf(md5Base62Format, i.Text(62))
}

// Md5Base62Upper returns a Base62 encoded MD5 hash of a string.
// Note Base62 encoding uses the GMP alphabet [0-9A-Za-z] instead
// of the Golang alphabet [0-9a-zA-Z] because the GMP alphabet
// may be more standard, e.g. used in GMP and follows ASCII
// table order.
func Md5Base62UpperFirst(s string) string {
	return stringsutil.ToOpposite(Md5Base62(s))
}
