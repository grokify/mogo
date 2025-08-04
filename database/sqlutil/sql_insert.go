package sqlutil

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
)

func BuildSQLXInsertSQLNamedParams(tblName string, colNames []string) (string, error) {
	tblName = strings.TrimSpace(tblName)
	if tblName == "" {
		return "", errors.New("table name cannot be empty")
	}

	// clone to prevent unintended side-effects if caller doesn't anticipate
	// colNames being sorted.
	colNamesLocal := slices.Clone(colNames)
	sort.Strings(colNamesLocal)
	var colNamesVars []string
	for _, colName := range colNamesLocal {
		if colName == "" {
			return "", errors.New("column name cannot be empty")
		} else if !IsUnquotedIdentifier(colName) {
			return "", fmt.Errorf("column name (%s) is not a valid unquoted identifier", colName)
		} else {
			colNamesVars = append(colNamesVars, ":"+colName)
		}
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tblName, strings.Join(colNamesLocal, ","), strings.Join(colNamesVars, ",")), nil
}
