package osutil

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/grokify/mogo/encoding/jsonutil"
)

func WriteFileJSON(filepath string, data any, perm os.FileMode, prefix, indent string) error {
	bytes, err := jsonutil.MarshalSimple(data, prefix, indent)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, bytes, perm)
}

var (
	ErrWriterNotInitialized = errors.New("bufio.Writer not initialized")
	ErrFileNotInitialized   = errors.New("os.File not initialized")
)

type FileWriter struct {
	File   *os.File
	Writer *bufio.Writer
}

func NewFileWriter(path string) (FileWriter, error) {
	fw := FileWriter{}
	file, err := os.Create(path)
	if err != nil {
		return fw, err
	}

	fw.File = file
	fw.Writer = bufio.NewWriter(file)

	return fw, nil
}

func (fw *FileWriter) WriteBytes(b []byte) (int, error) {
	return fw.Writer.Write(b)
}

func (fw *FileWriter) WriteString(addLinefeed bool, s ...string) (int, error) {
	n := 0
	for _, si := range s {
		if addLinefeed {
			si += "\n"
		}
		ni, err := fw.Writer.WriteString(si)
		if err != nil {
			return n, err
		}
		n += ni
	}
	return n, nil
}

func (fw *FileWriter) WriteStringf(addLinefeed bool, format string, a ...any) (int, error) {
	lf := ""
	if addLinefeed {
		lf = "\n"
	}
	return fw.Writer.WriteString(fmt.Sprintf(format, a...) + lf)
}

func (fw *FileWriter) Close() error {
	err := fw.Writer.Flush()
	if err != nil {
		return err
	}
	return fw.File.Close()
}

func WriteFileReader(filename string, r io.Reader, perm os.FileMode) (err error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	_, err = io.Copy(f, r)
	return err
}

// ReadFileSecure reads the file at the specified path after validating it does
// not contain path traversal sequences (".."). This is the recommended function
// for library code that receives paths from callers.
//
// Returns ErrPathTraversal (wrapped) if the path contains "..".
//
// For CLI entry points where the user explicitly provides paths, use
// os.ReadFile directly with a //nolint:gosec comment instead.
func ReadFileSecure(path string) ([]byte, error) {
	cleanPath, err := CleanPathSecure(path)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(cleanPath)
}

// WriteFileSecure writes data to the specified path after validating it does
// not contain path traversal sequences (".."). This is the recommended function
// for library code that receives paths from callers.
//
// Returns ErrPathTraversal (wrapped) if the path contains "..".
//
// For CLI entry points where the user explicitly provides paths, use
// os.WriteFile directly with a //nolint:gosec comment instead.
func WriteFileSecure(path string, data []byte, perm os.FileMode) error {
	cleanPath, err := CleanPathSecure(path)
	if err != nil {
		return err
	}
	return os.WriteFile(cleanPath, data, perm)
}
