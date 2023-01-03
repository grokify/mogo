package osutil

import (
	"bufio"
	"io"
	"os"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/io/ioutil"
)

func WriteFileJSON(filepath string, data interface{}, perm os.FileMode, prefix, indent string) error {
	bytes, err := jsonutil.MarshalSimple(data, prefix, indent)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, bytes, perm)
}

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

func (f *FileWriter) Close() {
	f.Writer.Flush()
	f.File.Close()
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
