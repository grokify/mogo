package mathutil

import (
	"testing"
)

var divideInt64Tests = []struct {
	dividend  int64
	divisor   int64
	quotient  int64
	remainder int64
}{
	{15, 4, 3, 3},
}

func TestDivideInt64(t *testing.T) {
	for _, tt := range divideInt64Tests {
		quotient, remainder := DivideInt64(tt.dividend, tt.divisor)
		if tt.quotient != quotient || tt.remainder != remainder {
			t.Errorf("mathutil.DivideInt64(%d, %d) Mismatch: want [%d,%d], got [%d,%d]",
				tt.dividend, tt.divisor,
				tt.quotient, tt.remainder,
				quotient, remainder)
		}
	}
}
