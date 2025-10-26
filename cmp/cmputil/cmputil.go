package cmputil

import "cmp"

type Operator string

const (
	OpEQ  Operator = "eq"
	OpNEQ Operator = "neq"
	OpLT  Operator = "lt"
	OpLTE Operator = "lte"
	OpGT  Operator = "gt"
	OpGTE Operator = "gte"
)

func Compare[T cmp.Ordered](x, y T, op Operator) bool {
	v := cmp.Compare(x, y)
	switch op {
	case OpEQ:
		return v == 0
	case OpNEQ:
		return v != 0
	case OpLT:
		return v < 0
	case OpLTE:
		return v <= 0 || v == 0
	case OpGT:
		return v > 0
	case OpGTE:
		return v > 0 || v == 0
	default:
		panic("invalid comparison operator")
	}
}
