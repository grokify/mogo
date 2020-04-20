package sql

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/grokify/gotilla/encoding/csvutil"
	"github.com/grokify/gotilla/type/stringsutil"
)

var maxInsertLength = 18000

var rxSplitLines = regexp.MustCompile(`(\r\n|\r|\n)`)

func SplitTextLines(text string) []string {
	return rxSplitLines.Split(text, -1)
}

func ReadFileCSVToSQLs(sqlFormat, filename, sep string, hasHeader, trimSpace bool, col uint) ([]string, error) {
	table, err := csvutil.NewTableDataFileSimple(filename, sep, hasHeader, trimSpace)
	if err != nil {
		return []string{}, err
	}
	values := []string{}
	for _, row := range table.Records {
		if len(row) > int(col) {
			values = append(values, row[col])
		}
	}
	values = stringsutil.SliceCondenseSpace(values, true, true)
	sqls := BuildSQLsInStrings(sqlFormat, values)
	return sqls, nil
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
	sqls := BuildSQLsInStrings(sqlFormat, values)
	return sqls, nil
}

func BuildSQLsInStrings(sqlFormat string, values []string) []string {
	sqls := []string{}
	sqlIns := SliceToSQLs(values)
	for _, sqlIn := range sqlIns {
		sqls = append(sqls, fmt.Sprintf(sqlFormat, sqlIn))
	}
	return sqls
}

func SliceToSQL(slice []string) string {
	newSlice := []string{}
	for _, el := range slice {
		newSlice = append(newSlice, "'"+el+"'")
	}
	return strings.Join(newSlice, ",")
}

func SliceToSQLs(slice []string) []string {
	max := maxInsertLength
	strIdx := 0
	newSlicesWip := [][]string{}
	newSlicesWip = append(newSlicesWip, []string{})
	for _, el := range slice {
		newStr := "'" + el + "'"
		if (LenStringForSlice(newSlicesWip[strIdx], ",") + len(newStr)) > max {
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
