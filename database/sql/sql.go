package sql

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/grokify/gotilla/encoding/csvutil"
	"github.com/grokify/gotilla/type/stringsutil"
)

const (
	MaxSQLLengthSOQL = 99500
	MaxSQLLengthApex = 18000
)

var rxSplitLines = regexp.MustCompile(`(\r\n|\r|\n)`)

func SplitTextLines(text string) []string {
	return rxSplitLines.Split(text, -1)
}

func ReadFileCSVToSQLs(sqlFormat, filename, sep string, hasHeader, trimSpace bool, col uint) ([]string, []string, error) {
	values, err := csvutil.ReadCSVFileSingleColumnValuesString(
		filename, sep, hasHeader, trimSpace, col, true)
	if err != nil {
		return []string{}, []string{}, err
	}

	sqls := BuildSQLsInStrings(sqlFormat, values, MaxSQLLengthSOQL)
	return sqls, values, nil
}

func ReadFileCSVToSQLsSimple(filename, sqlFormat string, hasHeader bool) ([]string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return []string{}, err
	}
	lines := stringsutil.SplitTextLines(string(bytes))
	if len(lines) == 0 {
		return []string{}, nil
	}
	if hasHeader {
		lines = lines[1:]
	}
	if len(lines) == 0 {
		return []string{}, nil
	}
	values := stringsutil.SliceCondenseSpace(lines, true, true)
	sqls := BuildSQLsInStrings(sqlFormat, values, MaxSQLLengthSOQL)
	return sqls, nil
}

func BuildSQLsInStrings(sqlFormat string, values []string, maxInsertLength int) []string {
	sqls := []string{}
	sqlIns := SliceToSQLs(values, maxInsertLength)
	for _, sqlIn := range sqlIns {
		sqls = append(sqls, fmt.Sprintf(sqlFormat, sqlIn))
	}
	return sqls
}

func SliceToSQL(slice []string) string {
	newSlice := []string{}
	quote := "'"
	for _, el := range slice {
		newSlice = append(newSlice, quote+el+quote)
	}
	return strings.Join(newSlice, ",")
}

func SliceToSQLs(slice []string, maxInsertLength int) []string {
	if maxInsertLength == 0 {
		maxInsertLength = MaxSQLLengthSOQL
	}
	quote := "'"
	strIdx := 0
	newSlicesWip := [][]string{}
	newSlicesWip = append(newSlicesWip, []string{})
	for _, el := range slice {
		newStr := quote + el + quote
		if maxInsertLength > 0 &&
			(LenStringForSlice(newSlicesWip[strIdx], ",")+len(newStr)) > maxInsertLength {
			newSlicesWip = append(newSlicesWip, []string{})
			strIdx += 1
		}
		newSlicesWip[strIdx] = append(newSlicesWip[strIdx], newStr)
	}
	newSlices := []string{}
	for _, slice := range newSlicesWip {
		newSlices = append(newSlices, strings.Join(slice, ","))
	}
	return newSlices
}

func LenStringForSlice(slice []string, sep string) int {
	return len(strings.Join(slice, sep))
}
