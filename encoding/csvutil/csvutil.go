package csvutil

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/type/stringsutil"
)

func NewReader(rs io.ReadSeeker, comma rune) (*csv.Reader, error) {
	var csvReader *csv.Reader
	// remove UTF-8 BOM, csv.Reader.Read() will return error = "line 1, column 1: bare \" in non-quoted-field"
	bom := make([]byte, 3)
	n, err := rs.Read(bom)
	if err != nil {
		return nil, err
	}
	if n < 3 || (n >= 3 && (bom[0] != 0xef || bom[1] != 0xbb || bom[2] != 0xbf)) {
		_, err = rs.Seek(0, 0) // Not a BOM -- seek back to the beginning
		if err != nil {
			return nil, err
		}
	}
	csvReader = csv.NewReader(rs)
	csvReader.Comma = comma
	return csvReader, nil
}

// NewReaderFile will create a csv.Reader and optionally strip off the
// byte order mark (BOM) if requested. Close file reader with
// `defer f.Close()`.
func NewReaderFile(path string, comma rune) (*csv.Reader, *os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, file, err
	}
	c, err := NewReader(file, comma)
	return c, file, err
}

func TrimUTF8ByteOrderMarkString(s string) string {
	byteOrderMarkAsString := string('\uFEFF')
	if strings.HasPrefix(s, byteOrderMarkAsString) {
		return strings.TrimPrefix(s, byteOrderMarkAsString)
	}
	return s
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

func (ch *CSVHeader) RecordMatch(row []string, andFilter map[string]stringsutil.MatchInfo) (bool, error) {
	for colName, matchInfo := range andFilter {
		idx := ch.Index(colName)
		if idx >= len(row) {
			return false, nil
		}
		match, err := stringsutil.Match(row[idx], matchInfo)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

func (ch *CSVHeader) RecordToMSS(row []string) map[string]string {
	mss := map[string]string{}
	l := len(row)
	for i, key := range ch.Columns {
		if i < l {
			mss[key] = row[i]
		}
	}
	return mss
}

func FilterCSVFile(inPath, outPath string, inComma rune, andFilter map[string]stringsutil.MatchInfo) error {
	reader, inFile, err := NewReaderFile(inPath, inComma)
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
		reader, inFile, err := NewReaderFile(inPath, inComma)
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

func MergeFilterCSVFilesToJSONL(inPaths []string, outPath string, inComma rune, andFilter map[string]stringsutil.MatchInfo) error {
	outFh, err := os.Create(outPath)
	if err != nil {
		return err
	}

	for _, inPath := range inPaths {
		reader, inFile, err := NewReaderFile(inPath, inComma)
		if err != nil {
			return err
		}

		csvHeader := CSVHeader{}
		j := -1

		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			j += 1

			if j == 0 {
				csvHeader.Columns = line
				continue
			}
			match, err := csvHeader.RecordMatch(line, andFilter)
			if err != nil {
				return err
			}
			if !match {
				continue
			}

			mss := csvHeader.RecordToMSS(line)

			bytes, err := json.Marshal(mss)
			if err != nil {
				return err
			}
			_, err = outFh.Write(bytes)
			if err != nil {
				return err
			}
			_, err = outFh.Write([]byte("\n"))
			if err != nil {
				return err
			}
			err = outFh.Sync()
			if err != nil {
				errOut := outFh.Close()
				errIn := inFile.Close()
				return errorsutil.Join(false, err, errOut, errIn)
			}
		}
		err = inFile.Close()
		if err != nil {
			errOut := outFh.Close()
			return errorsutil.Join(false, err, errOut)
		}
	}
	return outFh.Close()
}
