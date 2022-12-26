package slicesutil

// Dedupe returns a string slice with duplicate values removed. First observance is kept.
func Dedupe[S ~[]E, E comparable](s S) S {
	deduped := []E{}
	seen := map[E]int{}
	for _, val := range s {
		if _, ok := seen[val]; ok {
			continue
		}
		seen[val] = 1
		deduped = append(deduped, val)
	}
	return deduped
}

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
