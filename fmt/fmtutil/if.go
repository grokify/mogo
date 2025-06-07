package fmtutil

import (
	"fmt"
	"io"
)

// FprintIf is like `fmt.Fprint()` but doesn't fail when `io.Writer` is `nil`.
func FprintIf(w io.Writer, a ...any) (n int, err error) {
	if w != nil {
		return fmt.Fprint(w, a...)
	} else {
		return
	}
}

// FprintfIf is like `fmt.Fprintf()` but doesn't fail when `io.Writer` is `nil`.
func FprintfIf(w io.Writer, format string, a ...any) (n int, err error) {
	if w != nil {
		return fmt.Fprintf(w, format, a...)
	} else {
		return
	}
}
