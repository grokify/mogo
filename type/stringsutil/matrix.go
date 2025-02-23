package stringsutil

import (
	"errors"
	"fmt"
)

// Matrix2DColRowIndex returns the row index where the string supplied is first
// encountered for a supplied column index.
func Matrix2DColRowIndex[C comparable](mat [][]C, colIdx int, s C) (int, error) {
	if colIdx < 0 {
		return -1, errors.New("col index cannot be negative")
	}
	for y, row := range mat {
		if colIdx >= len(row) {
			return -1, fmt.Errorf("col index out of range [%d] with length %d", colIdx, len(row))
		}
		if row[colIdx] == s {
			return y, nil
		}
	}
	return -1, nil
}
