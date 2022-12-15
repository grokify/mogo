package stringsutil

import "fmt"

// Matrix2DColRowIndex returns the row index where the string supplied is first
// encountered for a supplied column index.
func Matrix2DColRowIndex[C comparable](mat [][]C, colIdx uint, s C) (int, error) {
	xint := int(colIdx)
	for y, row := range mat {
		if xint >= len(row) {
			return -1, fmt.Errorf("index out of range [%d] with length %d", xint, len(row))
		}
		if row[xint] == s {
			return y, nil
		}
	}
	return -1, nil
}
