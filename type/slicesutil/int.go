package slicesutil

func SliceIntSum(s []int) int {
	sum := 0
	for _, si := range s {
		sum += si
	}
	return sum
}

func MatrixIntColSums(m [][]int) []int {
	sums := []int{}
	for _, r := range m {
		for ci, c := range r {
			for ci >= len(sums) {
				sums = append(sums, 0)
			}
			sums[ci] += c
		}
	}
	return sums
}
