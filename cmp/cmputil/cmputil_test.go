package cmputil

import (
	"testing"
)

func TestCompare(t *testing.T) {
	var compareTests = []struct {
		x    int
		y    int
		op   Operator
		want bool
	}{
		{1, 1, OpEQ, true},
		{1, 1, OpNEQ, false},
		{1, 1, OpLT, false},
		{1, 1, OpLTE, true},
		{1, 1, OpGT, false},
		{1, 1, OpGTE, true},
		{1, 2, OpEQ, false},
		{1, 2, OpNEQ, true},
		{1, 2, OpLT, true},
		{1, 2, OpLTE, true},
		{1, 2, OpGT, false},
		{1, 2, OpGTE, false},
		{1, 0, OpEQ, false},
		{1, 0, OpNEQ, true},
		{1, 0, OpLT, false},
		{1, 0, OpLTE, false},
		{1, 0, OpGT, true},
		{1, 0, OpGTE, true},
	}

	for _, tt := range compareTests {
		try := Compare[int](tt.x, tt.y, tt.op)
		if try != tt.want {
			t.Errorf("cmputil.Compare(%d , %d,  %s) Mismatch: want (%v), got (%v)", tt.x, tt.y, string(tt.op), tt.want, try)
		}
	}
}
