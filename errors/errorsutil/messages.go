package errorsutil

import "fmt"

// ErrIndexOutOfRange returns an `error` using the standard message for
// index out of range.
func ErrIndexOutOfRange(idx, len int) error {
	return fmt.Errorf("index out of range [%d] with length %d", idx, len)
}
