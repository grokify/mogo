package ioutil

import (
	"bufio"
	"bytes"
	"io"
)

// ReaderToBytes reads from an io.Reader, e.g. io.ReadCloser
func ReaderToBytes(r io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// ReadAllOrError will successfully return the data
// or return the error in the value return value.
// This is useful to simply test scripts where the
// data is printed for debugging or testing.
func ReadAllOrError(r io.Reader) []byte {
	data, err := io.ReadAll(r)
	if err != nil {
		return []byte(err.Error())
	}
	return data
}

// Write writes from `Writer` to a `Reader`. See `osutil.WriteFileReader()`.
func Write(w *bufio.Writer, r io.Reader) error {
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		// write a chunk
		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
	}
	return w.Flush()
}
