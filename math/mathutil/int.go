package mathutil

import (
	"fmt"
)

// IntLen returns the string length of an integer.
func IntLen(i int) int {
	return len(fmt.Sprintf("%d", i))
}
