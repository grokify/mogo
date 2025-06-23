package mathutil

import (
	"fmt"
)

const (
	MaxInt63 = 2 ^ 62 - 1 // 4,611,686,018,427,387,903 . Since int63 is a signed integer, one bit is used for the sign. That leaves 62 bits for the value.
)

// IntLen returns the string length of an integer.
func IntLen(i int) int {
	return len(fmt.Sprintf("%d", i))
}
