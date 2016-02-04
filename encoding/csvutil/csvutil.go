package csvutil

import (
	"encoding/csv"
	"os"
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
