// md5util supports MD5 hashes in various formats.
package md5util

import (
	"crypto/md5" // #nosec G501
	"fmt"

	"github.com/grokify/mogo/math/bigint"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	// md5Base10Length int    = 39
	md5Base10Format string = `%039s`
	// Md5Base36Length is the length for a MD5 Base36 string
	// md5Base36Length int    = 25
	md5Base36Format string = `%025s`
	// md5Base62Length int    = 22
	md5Base62Format string = `%022s`
)

// MD5Base10  returns a Base10 encoded MD5 hash of a string.
func MD5Base10(s string) string {
	sum := md5.Sum([]byte(s))
	return fmt.Sprintf(md5Base10Format, bigint.MustEncodeToString(10, sum[:]))
	// i := new(big.Int)
	// i.SetString(fmt.Sprintf("%x", cryptomd5.Sum([]byte(s))), 16) // #nosec G401
	// return fmt.Sprintf(md5Base10Format, i.String())
}

// MD5Base36  returns a Base36 encoded MD5 hash of a string.
func MD5Base36(s string) string {
	sum := md5.Sum([]byte(s))
	return fmt.Sprintf(md5Base36Format, bigint.MustEncodeToString(36, sum[:]))
	// i := new(big.Int)
	// i.SetString(fmt.Sprintf("%x", cryptomd5.Sum([]byte(s))), 16) // #nosec G401
	// return fmt.Sprintf(md5Base36Format, i.Text(36))
}

// MD5Base62  returns a Base62 encoded MD5 hash of a string.
// This uses the Golang alphabet [0-9a-zA-Z].
func MD5Base62(s string) string {
	sum := md5.Sum([]byte(s))
	return fmt.Sprintf(md5Base62Format, bigint.MustEncodeToString(62, sum[:]))
	// i := new(big.Int)
	// i.SetString(fmt.Sprintf("%x", cryptomd5.Sum([]byte(s))), 16) // #nosec G401
	// return fmt.Sprintf(md5Base62Format, i.Text(62))
}

// MD5Base62UpperFirst returns a Base62 encoded MD5 hash of a string.
// Note Base62 encoding uses the GMP alphabet [0-9A-Za-z] instead
// of the Golang alphabet [0-9a-zA-Z] because the GMP alphabet
// may be more standard, e.g. used in GMP and follows ASCII
// table order.
func MD5Base62UpperFirst(s string) string {
	return stringsutil.ToOpposite(MD5Base62(s))
}
