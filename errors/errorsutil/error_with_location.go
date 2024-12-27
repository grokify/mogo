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

func NewErrorWithLocation(msg string) error {
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
