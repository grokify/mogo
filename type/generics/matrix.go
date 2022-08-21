package generics

// MakeMatrix2D returns a 2-dimensional matrix. Usage as follows
// `a := Make2D[uint8](dy, dx)`.
// Sourced from: https://stackoverflow.com/a/71781206/1908967
func MakeMatrix2D[T any](n, m int) [][]T {
	matrix := make([][]T, n)
	rows := make([]T, n*m)
	for i, startRow := 0, 0; i < n; i, startRow = i+1, startRow+m {
		endRow := startRow + m
		matrix[i] = rows[startRow:endRow:endRow]
	}
	return matrix
}
