package csvutil

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/grokify/gotilla/type/stringsutil"
)

/*

For UTF-8 BOM, csv.Reader.Read() will return error = "line 1, column 1: bare \" in non-quoted-field"

If you encounter this close the file and call again with stripBom = true

*/

// NewReader will create a csv.Reader and optionally strip off the
// byte order mark (BOM) if requested. Close file reader with
// `defer f.Close()`.
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

type CSVHeader struct {
	Columns []string
}

func (ch *CSVHeader) Index(want string) int {
	for i, try := range ch.Columns {
		if want == try {
			return i
		}
	}
	return -1
}

func (ch *CSVHeader) RowMatch(row []string, andFilter map[string]stringsutil.MatchInfo) bool {
	for colName, matchInfo := range andFilter {
		idx := ch.Index(colName)
		if idx >= len(row) {
			return false
		}
		if !stringsutil.Match(row[idx], matchInfo) {
			return false
		}
	}
	return true
}

func FilterCSVFile(inPath, outPath string, inComma rune, inStripBom bool, andFilter map[string]stringsutil.MatchInfo) error {
	reader, inFile, err := NewReader(inPath, inComma, inStripBom)
	if err != nil {
		return err
	}
	defer inFile.Close()
	writer, outFile, err := NewWriterFile(outPath)
	if err != nil {
		return err
	}
	defer writer.Flush()
	defer outFile.Close()
	return WriteCSVFiltered(reader, writer, andFilter, true)
}

// MergeFilterCSVFiles can merge and filter multiple CSV files. It expects row definitions to be the same
// across all input files.
func MergeFilterCSVFiles(inPaths []string, outPath string, inComma rune, inStripBom bool, andFilter map[string]stringsutil.MatchInfo) error {
	writer, outFile, err := NewWriterFile(outPath)
	if err != nil {
		return err
	}
	defer writer.Flush()
	defer outFile.Close()

	for i, inPath := range inPaths {
		reader, inFile, err := NewReader(inPath, inComma, inStripBom)
		if err != nil {
			return err
		}
		defer inFile.Close()

		writeHeader := false
		if i == 0 {
			writeHeader = true
		}
		err = WriteCSVFiltered(reader, writer, andFilter, writeHeader)
		if err != nil {
			return err
		}
	}
	return nil
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
		if csvHeader.RowMatch(line, andFilter) {
			err := writer.Write(line)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
