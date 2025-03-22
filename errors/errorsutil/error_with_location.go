package errorsutil

import (
	"fmt"
	"runtime"
)

// ErrorWithLocation represents an error message with code file and line locaiton.
// It is automatically populated when instantiated with `NewErrorWithLocation()`,
// and follows the `errors.Error` interface.
type ErrorWithLocation struct {
	Msg  string
	File string
	Line int
}

func (e *ErrorWithLocation) Error() string {
	return fmt.Sprintf("%s (at %s:%d)", e.Msg, e.File, e.Line)
}

func WrapWithLocation(err error) error {
	if err == nil {
		return nil
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = -1
	}
	return Wrap(err, fmt.Sprintf("error_location: file (%s) line (%d)", file, line))
}

func NewWithLocation(msg string) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = -1
	}
	return &ErrorWithLocation{
		Msg:  msg,
		File: file,
		Line: line,
	}
}
