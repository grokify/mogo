package csvutil

import (
	"encoding/csv"
	"os"
)

func NewReader(path string, comma rune) (*csv.Reader, *os.File, error) {
	var myCsv *csv.Reader
	var file *os.File
	file, err := os.Open(path)
	if err != nil {
		return myCsv, file, err
	}
	reader := csv.NewReader(file)
	reader.Comma = comma
	return reader, file, nil
}
