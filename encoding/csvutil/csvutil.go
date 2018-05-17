package csvutil

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/grokify/gotilla/strings/stringsutil"
)

/*

For UTF-8 BOM, csv.Reader.Read() will return error = "line 1, column 1: bare \" in non-quoted-field"

If you encounter this close the file and call again with stripBom = true

*/

// NewReader will create a csv.Reader and optionally strip off the
// byte order mark (BOM) if requested.
func NewReader(path string, comma rune, stripBom bool) (*csv.Reader, *os.File, error) {
	var csvReader *csv.Reader
	var file *os.File
	file, err := os.Open(path)
	if err != nil {
		return csvReader, file, err
	}
	if stripBom {
		b3 := make([]byte, 3)
		_, err := file.Read(b3)
		if err != nil {
			return csvReader, file, err
		}
	}
	csvReader = csv.NewReader(file)
	csvReader.Comma = comma
	return csvReader, file, nil
}

// Writer is a struct for a CSV/TSV writer.
type Writer struct {
	Separator        string
	StripRepeatedSep bool
	ReplaceSeparator bool
	SeparatorAlt     string
	File             *os.File
}

// NewWriter returns a Writer with the separator params set.
func NewWriter(filepath, sep string, replaceSeparator bool, alt string) (Writer, error) {
	w := Writer{
		Separator:        sep,
		StripRepeatedSep: false,
		ReplaceSeparator: replaceSeparator,
		SeparatorAlt:     alt}
	err := w.open(filepath)
	return w, err
}

// Open opens a filepath.
func (w *Writer) open(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	w.File = f
	return nil
}

// AddLine adds an []interface{} to the file.
func (w *Writer) AddLine(cells []interface{}) {
	fmt.Fprintf(
		w.File,
		"%s\n",
		stringsutil.JoinInterface(cells, w.Separator, false, w.ReplaceSeparator, w.SeparatorAlt))
}

// Close closes the file.
func (w *Writer) Close() {
	w.File.Close()
}
