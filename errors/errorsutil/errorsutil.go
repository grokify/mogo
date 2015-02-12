package errorsutil

import (
	"errors"
	"fmt"
)

func Append(err error, str string) error {
	return errors.New(fmt.Sprint(err) + str)
}
