package csvutil

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/grokify/simplego/type/stringsutil"
)

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
	return w, w.open(filepath)
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
func (w *Writer) AddLine(cells []interface{}) error {
	_, err := w.File.WriteString(
		stringsutil.JoinInterface(
			cells, w.Separator, false, w.ReplaceSeparator, w.SeparatorAlt) + "\n")
	if err != nil {
		return err
	}
	return w.File.Sync()
}

// Close closes the file.
func (w *Writer) Close() error {
	return w.File.Close()
}

func NewWriterFile(filename string) (*csv.Writer, *os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, file, err
	}
	//defer file.Close()

	writer := csv.NewWriter(file)
	//defer writer.Flush()
	return writer, file, nil
}

// WriteCSVFiltered filters an existing CSV and writes the matching lines
// to a *csv.Writer.
func WriteCSVFiltered(reader *csv.Reader, writer *csv.Writer, andFilter map[string]stringsutil.MatchInfo, writeHeader bool) error {
	csvHeader := CSVHeader{}
	i := -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		i += 1
		if i == 0 {
			csvHeader.Columns = line
			if writeHeader {
				err := writer.Write(line)
				if err != nil {
					return err
				}
			}
			continue
		}
		match, err := csvHeader.RecordMatch(line, andFilter)
		if err != nil {
			return err
		}
		if match {
			err := writer.Write(line)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
