package number

func MatrixRowsMax(d [][]float64) []float64 {
	var rows []float64
	if len(d) == 0 {
		return rows
	}
	for y := 0; y < len(d); y++ {
		var rowDistMax float64
		init := false
		for x := 0; x < len(d[0]); x++ {
			if !init {
				rowDistMax = d[y][x]
				init = true
			} else if d[y][x] > rowDistMax {
				rowDistMax = d[y][x]
			}
		}
		rows = append(rows, rowDistMax)
	}
	return rows
}

func MatrixColsMax(d [][]float64) []float64 {
	var cols []float64
	if len(d) == 0 {
		return cols
	} else if len(d[0]) == 0 {
		return cols
	}
	for x := 0; x < len(d[0]); x++ {
		var colDistMax float64
		init := false
		for y := 0; y < len(d); y++ {
			if !init {
				colDistMax = d[y][x]
				init = true
			} else if d[y][x] > colDistMax {
				colDistMax = d[y][x]
			}
		}
		cols = append(cols, colDistMax)
	}
	return cols
}
