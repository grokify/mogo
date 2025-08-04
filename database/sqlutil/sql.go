package sqlutil

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/grokify/mogo/encoding/csvutil"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	MaxSQLLengthSOQL = 99500
	MaxSQLLengthApex = 18000
)

/*
var rxSplitLines = regexp.MustCompile(`(\r\n|\r|\n)`)

func SplitTextLines(text string) []string {
	return rxSplitLines.Split(text, -1)
}
*/

/*

func ReadCSVFileSingleColumnValuesString(filename string, sep rune, stripBOM, hasHeader, trimSpace bool, colIdx uint, condenseUniqueSort bool) ([]string, error) {

*/

func ReadFileCSVToSQLs(sqlFormat, filename string, sep rune, hasHeader, trimSpace bool, colIdx uint32) ([]string, []string, error) {
	values, err := csvutil.ReadCSVFileSingleColumnValuesString(
		filename, sep, hasHeader, trimSpace, colIdx, true)
	if err != nil {
		return []string{}, []string{}, err
	}

	sqls := BuildSQLsInStrings(sqlFormat, values, MaxSQLLengthSOQL)
	return sqls, values, nil
}

func ReadFileCSVToSQLsSimple(filename, sqlFormat string, hasHeader bool) ([]string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return []string{}, err
	}
	lines := stringsutil.SplitLines(string(bytes))
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

// BuildSQLsInStrings returns a SQL statement inserting SQL IN values.
// The supplied SQL fromat should have an IN clause like `IN (%s)`.
// For example `SELECT Id from Account WHERE Id IN (%s)`.
func BuildSQLsInStrings(sqlFormat string, values []string, maxInsertLength int) []string {
	sqls := []string{}
	sqlIns := SliceToSQLs(values, maxInsertLength)
	for _, sqlIn := range sqlIns {
		sqls = append(sqls, fmt.Sprintf(sqlFormat, sqlIn))
	}
	return sqls
}

const SingleQuote = "'"

var rxSingleQuote = regexp.MustCompile(`'`)

func QuoteWord(s string) string {
	return SingleQuote +
		rxSingleQuote.ReplaceAllString(s, "''") +
		SingleQuote
}

func SliceToSQL(slice []string) string {
	newSlice := []string{}
	for _, el := range slice {
		newSlice = append(newSlice, QuoteWord(el))
	}
	return strings.Join(newSlice, ",")
}

// SliceToSQLs returns a slice of string elements separated by
// commas without the surrounding parentheses `()`. This is useful for breaking up
// a long SQL IN clause into multiple SQL statements.
func SliceToSQLs(slice []string, maxInsertLength int) []string {
	if maxInsertLength <= 0 {
		maxInsertLength = MaxSQLLengthSOQL
	}
	strIdx := 0
	newSlicesWip := [][]string{}
	newSlicesWip = append(newSlicesWip, []string{})
	for _, el := range slice {
		newStr := QuoteWord(el)
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
