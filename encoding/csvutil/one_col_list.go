package csvutil

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

// ReadFileOneColListToGrid parses a file with one value per row.
func ReadFileOneColListToGrid(filename string, colCount int, validateLength bool) ([][]string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return [][]string{}, err
	}
	return ParseOneColListToGrid(strings.Split(string(bytes), "\n"), colCount, validateLength)
}

// ParseOneColListToGrid parses a set of lines with one value per row.
func ParseOneColListToGrid(lines []string, colCount int, validateLength bool) ([][]string, error) {
	rows := [][]string{}
	curRow := []string{}
	colCountFloat := float64(colCount)

	for i, line := range lines {
		curRow = append(curRow, line)
		if math.Mod(float64(i+1), colCountFloat) == 0 {
			if len(curRow) > 0 {
				rows = append(rows, curRow)
				curRow = []string{}
			}
		}
	}
	if len(curRow) > 0 {
		rows = append(rows, curRow)
		curRow = []string{}
	}
	if validateLength {
		for _, row := range rows {
			if len(row) != colCount {
				return rows, fmt.Errorf("Error Row Length Want [%d] Have [%d]", colCount, len(row))
			}
		}
	}
	return rows, nil
}
