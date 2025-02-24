package ioutil

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
)

type AtReader interface {
	io.Reader
	io.ReaderAt
}

func IsReader(i any) bool {
	reader := reflect.TypeOf((*io.Reader)(nil)).Elem()
	return reflect.PointerTo(reflect.TypeOf(i).Elem()).Implements(reader)
}

func MustPrintReader(r io.Reader) {
	if b, err := io.ReadAll(r); err != nil {
		panic(err)
	} else {
		fmt.Println(string(b))
	}
}

/*
// ReaderToBytes reads from an io.Reader, e.g. io.ReadCloser
func ReaderToBytes(r io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}
*/

// ReadAllOrError will successfully return the data
// or return the error in the value return value.
// This is useful to simply test scripts where the
// data is printed for debugging or testing.
func ReadAllOrError(r io.Reader) []byte {
	if b, err := io.ReadAll(r); err != nil {
		return []byte(err.Error())
	} else {
		return b
	}
}

// ReadLimit returns the first `limit` bytes from a reader.
func ReadLimit(r io.Reader, limit int64) ([]byte, error) {
	if limit <= 0 {
		return []byte{}, nil
	}
	return io.ReadAll(io.LimitReader(r, limit))
}

// ReaderToReadSeeker converts an `io.Reader` to an `io.ReadSeeker`. It does this
// by reading all data in `io.Reader`.
func ReaderToReadSeeker(r io.Reader) (io.ReadSeeker, error) {
	if b, err := io.ReadAll(r); err != nil {
		return nil, err
	} else {
		return bytes.NewReader(b), nil
	}
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

// SkipWriter is an `io.Writer` that skips writing the first `n` bytes passed. This is
// useful when the `io.Writer` writes some undesirable data which will be omitted with
// this functionality.
type SkipWriter struct {
	// Rewritten from the following under MIT license: https://github.com/jdeng/goheif/blob/a0d6a8b3e68f9d613abd9ae1db63c72ba33abd14/heic2jpg/main.go
	w      io.Writer
	offset int
}

// NewSkipWriter returns an SkipWriter that writes to w skipping the first offset off bytes.
func NewSkipWriter(w io.Writer, off int) *SkipWriter {
	return &SkipWriter{w: w, offset: off}
}

// Write fulfills the `io.Writer` interface.
func (s *SkipWriter) Write(p []byte) (n int, err error) {
	if s.offset <= 0 {
		n, err = s.w.Write(p)
		return
	}

	if plen := len(p); plen < s.offset {
		s.offset -= plen
		n = plen
		return
	}

	n, err = s.w.Write(p[s.offset:])
	if err != nil {
		return
	}
	n += s.offset
	s.offset = 0
	return
}
