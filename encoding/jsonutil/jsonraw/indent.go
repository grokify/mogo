package jsonraw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
)

// IndentBytes converts a JSON byte array into a prettified byte array.
func IndentBytes(data []byte, prefix, indent string) ([]byte, error) {
	var out bytes.Buffer
	if len(data) == 0 {
		return []byte{}, nil
	} else if err := json.Indent(&out, data, prefix, indent); err != nil {
		return []byte{}, err
	} else {
		return out.Bytes(), nil
	}
}

func MarshalFileIndentBytes(name string, data []byte, prefix, indent string, perm fs.FileMode) error {
	if data, err := IndentBytes(data, prefix, indent); err != nil {
		return err
	} else {
		return os.WriteFile(name, data, perm)
	}
}

// Indent returns a byte slice of indented JSON given an `io.Reader`.
// It is useful to use with `http.Response.Body` which is an `io.ReadCloser`.
func Indent(r io.Reader, prefix, indent string) ([]byte, error) {
	if b, err := io.ReadAll(r); err != nil {
		return b, err
	} else {
		return IndentBytes(b, prefix, indent)
	}
}

// PrintIndent returns an indented JSON byte array given an `io.Reader`.
func PrintIndent(r io.Reader, prefix, indent string) ([]byte, error) {
	outBytes, err := Indent(r, prefix, indent)
	if err != nil {
		return []byte{}, err
	}
	_, err = fmt.Println(string(outBytes))
	return outBytes, err
}
