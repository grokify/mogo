package number

import (
	"errors"

	"github.com/grokify/mogo/strconv/strconvutil"
	"golang.org/x/exp/constraints"
)

type Integers[C constraints.Integer] []C

func (s Integers[C]) Sum() C {
	var sum C
	for _, v := range s {
		sum += v
	}
	return sum
}

func (s Integers[C]) SumString() string {
	return strconvutil.Itoa(s.Sum())
}

type IntegersMatrix[C constraints.Integer] [][]C

// ColumnSums returns a slice where the element in each value is the sum of the
// column, e.g. [{sumColumn0},{sumColumn1},{sumColumn2}]
func (m IntegersMatrix[C]) ColumnSums(colIdx int, zeroOnEmpty bool) ([]C, error) {
	if colIdx < 0 {
		if zeroOnEmpty {
			return []C{}, nil
		} else {
			return []C{}, errors.New("column index cannot be zero")
		}
	}
	var sums []C
	var zero C
	for _, r := range m {
		if len(r) == 0 || len(r)-1 > colIdx {
			if zeroOnEmpty {
				sums = append(sums, zero)
				continue
			} else {
				return []C{}, errors.New("column index for row is missing")
			}
		}
		for i, ri := range r {
			if i < len(sums) {
				sums[i] += ri
			} else if i == len(sums) {
				sums = append(sums, ri)
			} else {
				panic("internal error")
			}
		}
	}
	return sums, nil
}

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
