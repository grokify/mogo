package errorsutil

import (
	"errors"
	"fmt"
)

// Append adds additional text to an existing error.
func Append(err error, str string) error {
	return errors.New(fmt.Sprint(err) + str)
}
