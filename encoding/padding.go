// encoding provides generic encoding support.
package encoding

import "github.com/grokify/mogo/math/mathutil"

func Pad4(encoded, char string) string {
	inputLength := len(encoded)
	_, rem := mathutil.Divide(int64(inputLength), int64(4))
	if rem == 0 {
		return encoded
	}
	if len(char) == 0 {
		char = " "
	}
	switch rem {
	case 1:
		rem = 3
	case 3:
		rem = 1
	}
	for i := 0; i < int(rem); i++ {
		encoded += char
	}
	return encoded[:inputLength+int(rem)]
}
