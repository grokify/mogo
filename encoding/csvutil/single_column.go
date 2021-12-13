package csvutil

import (
	"fmt"
	"io"
	"strings"

	"github.com/grokify/mogo/io/ioutilmore"
	"github.com/grokify/mogo/type/stringsutil"
)

func ReadCSVFileSingleColumnValuesString(filename string, sep rune, hasHeader, trimSpace bool, colIdx uint, condenseUniqueSort bool) ([]string, error) {
	values := []string{}
	csvReader, file, err := NewReaderFile(filename, sep)
	if err != nil {
		return values, nil
	}

	i := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return values, ioutilmore.CloseFileWithError(file, err)
		}
		if hasHeader && i == 0 {
			i++
			continue
		}

		if int(colIdx) >= len(record) {
			return values, fmt.Errorf("E_RECORD_TOO_SHORT LEN[%v] WANT_INDEX [%v]", len(record), colIdx)
		}
		val := record[colIdx]
		if trimSpace {
			val = strings.TrimSpace(val)
		}
		values = append(values, val)
	}
	/*
		tbl, err := NewTableDataFileSimple(filename, sep, hasHeader, trimSpace)
		if err != nil {
			return []string{}, err
		}

		for _, row := range tbl.Records {
			if len(row) > int(col) {
				values = append(values, row[col])
			}
		}*/
	if condenseUniqueSort {
		values = stringsutil.SliceCondenseSpace(values, true, true)
	}
	return values, nil
}

func ReadCSVFilesSingleColumnValuesString(files []string, sep rune, hasHeader, trimSpace bool, colIdx uint, condenseUniqueSort bool) ([]string, error) {
	values := []string{}
	for _, file := range files {
		fileValues, err := ReadCSVFileSingleColumnValuesString(
			file, sep, hasHeader, trimSpace, colIdx, false)
		if err != nil {
			return values, err
		}
		values = append(values, fileValues...)
	}
	if condenseUniqueSort {
		values = stringsutil.SliceCondenseSpace(values, true, true)
	}
	return values, nil
}
