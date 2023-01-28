package osutil

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/io/ioutil"
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

func WriteFileReader(filename string, r io.Reader) error {
	// https://stackoverflow.com/questions/1821811/how-to-read-write-from-to-a-file-using-go
	// open output file
	// fo, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, perm)
	fo, err := os.Create(filename)
	if err != nil {
		return err
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	w := bufio.NewWriter(fo)
	return ioutil.Write(w, r)
}
