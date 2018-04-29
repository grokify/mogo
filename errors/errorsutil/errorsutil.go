package errorsutil

import (
	"errors"
	"fmt"
)

// Append adds additional text to an existing error.
func Append(err error, str string) error {
	return errors.New(fmt.Sprint(err) + str)
}

// PanicOnErr is a syntactic sugar function to panic on error.
func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

type ErrorInfo struct {
	Error   error
	Code    string
	Display string
	Input   string
	Correct string
}
